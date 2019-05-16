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
	docMeta *documentMetadata
	errors  <-chan error
	found   <-chan *URI
}
