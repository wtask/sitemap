package sitemap

import (
	"time"
)

type task struct {
	*URI
	level int
}

type documentMetadata struct {
	modified *time.Time
}

type taskResult struct {
	task
	err error
	meta     *documentMetadata // task document metadata
	// errors <-chan error
	found <-chan *task
}

// func worker(t task, depth uint) taskResult {
// 	found := make(chan *task)
// }
