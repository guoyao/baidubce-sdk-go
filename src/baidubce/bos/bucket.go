package bos

import (
	"time"
)

type Location struct {
	LocationConstraint string
}

type BucketSummary struct {
	Owner   map[string]string
	Buckets []Bucket
}

func (bucketSummary *BucketSummary) GetBuckets() []Bucket {
	return bucketSummary.Buckets
}

type Bucket struct {
	Name, Location string
	CreationDate   time.Time
}
