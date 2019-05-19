package sitemap

import "time"

// Target - represents a link that corresponds to a specific level of hierarchy
type Target struct {
	*URI
	Level uint
}

// DocumentMeta - metadata of fetched targets
type DocumentMeta struct {
	// Modified - document modification time
	Modified time.Time
}

// completedTarget - processed target data
type completedTarget struct {
	Target
	err  error
	meta *DocumentMeta // task document metadata
	// errors <-chan error
	targets <-chan Target
}

// MapItem - final result of site map
type MapItem struct {
	Target
	*DocumentMeta
}
