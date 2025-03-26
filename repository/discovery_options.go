package repository

import (
	"repo-file-sync/set"
)

type DiscoveryOptions struct {
	IncludePatterns set.Set[string]
}

func NewDiscoveryOptions(includePatterns []string) DiscoveryOptions {
	return DiscoveryOptions{
		IncludePatterns: set.FromSlice(includePatterns),
	}
}

func DefaultGlobalDiscoveryOptions() DiscoveryOptions {
	return DiscoveryOptions{
		IncludePatterns: set.FromSlice(defaultIncludePatterns),
	}
}

func (do *DiscoveryOptions) Merge(other DiscoveryOptions) DiscoveryOptions {
	return DiscoveryOptions{
		IncludePatterns: other.IncludePatterns.Copy().Merge(do.IncludePatterns),
	}
}
