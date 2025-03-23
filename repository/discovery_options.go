package repository

import (
	"repo-file-sync/set"
)

type DiscoveryOptions struct {
	IncludePatterns set.Set[string]
	ExcludePatterns set.Set[string]
}

func NewDiscoveryOptions(includePatterns []string) DiscoveryOptions {
	return DiscoveryOptions{
		IncludePatterns: set.FromSlice(includePatterns),
		ExcludePatterns: set.New[string](),
	}
}

func DefaultGlobalDiscoveryOptions() DiscoveryOptions {
	return DiscoveryOptions{
		IncludePatterns: set.FromSlice(defaultIncludePatterns),
		ExcludePatterns: set.FromSlice(defaultExcludePatterns),
	}
}

func (do *DiscoveryOptions) Merge(other DiscoveryOptions) DiscoveryOptions {
	return DiscoveryOptions{
		IncludePatterns: other.IncludePatterns.Copy().Merge(do.IncludePatterns),
		ExcludePatterns: other.ExcludePatterns.Copy().Merge(do.ExcludePatterns),
	}
}
