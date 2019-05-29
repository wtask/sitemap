package sitemap

import (
	"fmt"
	"net/url"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Parser - represent type to explore and build site map.
type Parser struct {
	errorHandler   ErrorHandler  // async
	requestTimeout time.Duration // optional
	queueCap       uint
}

const (
	// DefaultNumWorkers - default num of goroutines which fetches and parses html documents.
	DefaultNumWorkers uint = 4
	// DefaultQueueCap - default capacity of internal queue.
	DefaultQueueCap uint = 1000
)

// ErrorHandler - func which should handle parsing error.
type ErrorHandler func(error)

type parserOption func(*Parser) error

func failedOption(err error) parserOption {
	return func(p *Parser) error {
		return err
	}
}

func (p *Parser) setup(options ...parserOption) error {
	if p == nil {
		return nil
	}
	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(p); err != nil {
			return err
		}
	}
	return nil
}

// WithErrorHandler - specify parsing error handler.
// By default, all occurred errors are ignored.
// Note: error handler will start inside own goroutine, which will silently recovered in case of panic.
// Also the Parse method will wait for handlers are completed before returning results.
func WithErrorHandler(h ErrorHandler) parserOption {
	return func(p *Parser) error {
		p.errorHandler = h
		return nil
	}
}

// WithRequestTimeout - declare timeout for any requests made with Parser.
func WithRequestTimeout(timeout time.Duration) parserOption {
	if timeout < 0 {
		return failedOption(fmt.Errorf("Invalid request timeout %v", timeout))
	}
	return func(p *Parser) error {
		p.requestTimeout = timeout
		return nil
	}
}

// NewParser - create Parser instance with optional features.
func NewParser(options ...parserOption) (*Parser, error) {
	p := &Parser{queueCap: DefaultQueueCap}
	if err := p.setup(options...); err != nil {
		return nil, err
	}
	return p, nil
}

// Parse - takes root URI and max depth to find all links inside html documents available from root.
// TODO Add support to cancel parsing with help of context.Context
func (p *Parser) Parse(root *URI, depth, workers uint) []MapItem {
	// TODO resolve p == nil case
	if workers == 0 {
		// or panic?
		workers = DefaultNumWorkers
	}
	if p.queueCap == 0 {
		p.queueCap = DefaultQueueCap
	}
	queue := make(chan Target, p.queueCap)
	pending := make(chan (<-chan Target), workers)
	results := sync.Map{}

	num := struct{ workers, fillers int64 }{0, 0} // goroutines counters
	eh := sync.WaitGroup{}

	ensureWorkers := func() {
		for i := atomic.LoadInt64(&num.workers); i < int64(workers); i++ {
			select {
			case target := <-queue:
				atomic.AddInt64(&num.workers, 1)
				go func() {
					defer func() {
						atomic.AddInt64(&num.workers, -1)
					}()
					if _, ok := results.Load(target.URI.String()); ok {
						// already have target, drop it
						return
					}
					completed := p.worker(root, depth, target)
					if completed.err != nil && p.errorHandler != nil {
						eh.Add(1)
						go func() {
							defer func() {
								recover() // protect parser from handler panic
								eh.Done()
							}()
							p.errorHandler(completed.err)
						}()
					}
					// Do not check "loaded" result here,
					// always send completed.targets into pending chan to prevent leak of background goroutines.
					results.LoadOrStore(
						completed.Target.URI.String(),
						MapItem{completed.Target.URI, completed.meta},
					)
					if completed.targets != nil {
						pending <- completed.targets
					}
				}()
			default:
				return
			}
		}
	}

	ensureFillers := func() {
		select {
		// fill the queue when there are pending targets
		case targets := <-pending:
			atomic.AddInt64(&num.fillers, 1)
			go func() {
				defer func() {
					atomic.AddInt64(&num.fillers, -1)
				}()
				for t := range targets {
					queue <- t
				}
			}()
		default:
		}
	}

	queue <- Target{root, 0}
	for {
		// main loop
		ensureWorkers()
		ensureFillers()

		if len(queue) == 0 &&
			len(pending) == 0 &&
			atomic.LoadInt64(&num.workers) == 0 &&
			atomic.LoadInt64(&num.fillers) == 0 {
			// all done
			break
		}
	}

	found := []MapItem{}
	results.Range(func(_ interface{}, value interface{}) bool {
		if item, ok := value.(MapItem); ok {
			found = append(found, item)
		}
		return true
	})

	eh.Wait()

	return found
}

// worker - fetches and parses target document.
// Arguments `root` and `depth` are required to build absolute URI properly.
func (p *Parser) worker(root *URI, depth uint, t Target) completedTarget {
	// if an error occurred, the doc could still be partially exists,
	// below we will check doc body
	doc, meta, err := fetchDocument(t.URI, p.requestTimeout)

	result := completedTarget{
		Target:  t,
		err:     err,
		meta:    meta,
		targets: nil,
	}

	if t.Level >= depth {
		// stop parsing
		return result
	}

	targets := make(chan Target)
	result.targets = targets

	go func() {
		defer close(targets)
		body := firstNode("body", doc)
		if body == nil {
			return
		}
		base, _ := NewURI(
			attribute("href", firstNode("base", firstNode("head", doc))),
		)
		for _, href := range collectAttributes("a", "href", body, nil) {
			var url *url.URL
			if base != nil {
				url, _ = base.Parse(href)
			} else {
				url, _ = root.Parse(href)
			}
			if url == nil {
				continue
			}
			url.Fragment = "" // always drop fragment
			link, err := NewURI(url.String())
			if err != nil {
				continue
			}
			if root.Scheme != link.Scheme ||
				root.Hostname() != link.Hostname() ||
				!strings.HasPrefix(
					path.Dir(link.EscapedPath()),
					path.Dir(root.EscapedPath()),
				) {
				// TODO Make more reliable verification for nested targets
				// Add method to URI
				continue
			}
			targets <- Target{link, t.Level + 1}
		}
	}()

	return result
}
