package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	mapper("http:/localhost/", 1, 4)
}

func mapper(root string, depth, workerLimit int) {
	queue := make(chan Job, 1000)
	// pool := make(chan Job, workerLimit)
	results := make(chan Result, workerLimit)
	found := sync.Map{}

	queue <- Job{Link{url: root, level: 0}}

	pool := func() chan Job {
		// формируем пул для работы воркеров
		p := make(chan Job, workerLimit)
		go func() {
			defer close(p)
			for j := range queue {
				fmt.Printf("Send job from queue to pool: %s\n", j.Link.url)
				// TODO добавить select, чтобы остановиться, если <- quit
				p <- j
			}
			fmt.Println("Quit from pool filler")
		}()
		return p
	}()

	var goroutines int32
	wg := sync.WaitGroup{}

	for completed := false; !completed; {
		select {
		case j := <-pool:
			wg.Add(1)
			atomic.AddInt32(&goroutines, 1)
			go func(j Job, results chan<- Result) {
				defer func() {
					wg.Done()
					atomic.AddInt32(&goroutines, -1)
				}()
				results <- worker(j)
			}(j, results)
		case r := <-results:
			wg.Add(1)
			atomic.AddInt32(&goroutines, 1)
			go func() {
				defer func() {
					wg.Done()
					atomic.AddInt32(&goroutines, -1)
				}()
				fmt.Printf("Got result in background: %s\n", r.Job.Link.url)
				if _, ok := found.Load(r.Job.Link.url); !ok {
					fmt.Printf("Add result to map: %s\n", r.Job.Link.url)
					found.Store(r.Job.Link.url, true)
				}
				if r.err != nil {
					// TODO log job error
				}
				for l := range r.found {
					if r.Link.level > depth {
						fmt.Printf("Dropped link due to level is %d\n", r.Link.level)
						// дропаем ссылки, освобождая канал
						continue
					}
					if _, ok := found.Load(l.url); !ok {
						fmt.Printf("Queue link: %s\n", l.url)
						queue <- Job{Link: l}
					}
				}
			}()
		default:
			fmt.Println("No job, No results")
		}

		completed = len(queue) == 0 && len(pool) == 0 && len(results) == 0 && atomic.LoadInt32(&goroutines) == 0
	}

	wg.Wait()
	close(queue) // также закроет pool
	close(results)
}

type Link struct {
	url   string
	level int
	// TODO add last-modified and other attribbutes
}

type Job struct {
	// quit func() bool
	Link
}

type Result struct {
	// обработанный результат, который надо только сохранить
	Job
	err   error
	found <-chan Link // канал результатов генератора
}

func worker(j Job) Result {
	time.Sleep(1 * time.Second)
	found := make(chan Link) // TODO можно использовать буферизованный
	go func(total int) {
		defer close(found)
		for i := 0; i < total; i++ {
			found <- Link{
				url:   fmt.Sprintf("%s%d-%d/", j.Link.url, j.Link.level, i),
				level: j.Link.level + 1,
			}
		}
	}(5)
	return Result{
		Job:   j,
		err:   nil,
		found: found,
	}
}
