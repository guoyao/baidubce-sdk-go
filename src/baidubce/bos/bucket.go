package bos

import (
	"time"
)

// Location is a struct for bucket location info.
type Location struct {
	LocationConstraint string
}

type BucketOwner struct {
	Id          string
	DisplayName string
}

// Bucket is a struct for bucket info.
type Bucket struct {
	Name, Location string
	CreationDate   time.Time
}

// BucketSummary is a struct for bucket summary.
type BucketSummary struct {
	Owner   BucketOwner
	Buckets []Bucket
}
