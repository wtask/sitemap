package sitemap

import (
	"context"
	"fmt"
	"time"
)

// Parser - repsent type to explore and build site map.
type Parser struct {
	сtx            *context.Context // optional
	errors         chan<- error     // optional
	requestTimeout time.Duration    // optional
	queueLen       int
}

const defaultQueueLen = 1000

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

// WithContext - add optional context to Parser to allow cancellation.
func WithContext(ctx context.Context) parserOption {
	return func(p *Parser) error {
		p.сtx = &ctx
		return nil
	}
}

// WithErrorChannel - add optional channel to send there parsing errors.
func WithErrorChannel(errors chan<- error) parserOption {
	return func(p *Parser) error {
		p.errors = errors
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
	p := &Parser{queueLen: defaultQueueLen}
	if err := p.setup(options...); err != nil {
		return nil, err
	}
	return p, nil
}

func (p Parser) Parse(root *URI, depth, workers int) []*URI {
	return nil
}
