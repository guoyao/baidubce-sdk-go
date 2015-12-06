package bos

import (
	"time"
)

// Location is a struct for bucket location info.
type Location struct {
	LocationConstraint string
}

// BucketSummary is a struct for bucket summary.
type BucketSummary struct {
	Owner   map[string]string
	Buckets []Bucket
}

// Bucket is a struct for bucket info.
type Bucket struct {
	Name, Location string
	CreationDate   time.Time
}
