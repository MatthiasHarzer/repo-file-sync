package repository

import (
	"repo-file-sync/set"
)

type DiscoveryOptions struct {
	IncludePatterns set.Set[string]
	ExcludePatterns set.Set[string]
}

func NewDiscoveryOptions() DiscoveryOptions {
	return DiscoveryOptions{
		IncludePatterns: set.New[string](),
		ExcludePatterns: set.New[string](),
	}
}

func (do *DiscoveryOptions) Merge(other DiscoveryOptions) DiscoveryOptions {
	return DiscoveryOptions{
		IncludePatterns: other.IncludePatterns.Copy().Merge(do.IncludePatterns),
		ExcludePatterns: other.ExcludePatterns.Copy().Merge(do.ExcludePatterns),
	}
}
