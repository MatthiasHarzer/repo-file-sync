package repository

import (
	"ide-config-sync/set"
)

// TODO: Maybe there is a better name for this struct?
type Options struct {
	IncludePatterns set.Set[string]
	ExcludePatterns set.Set[string]
}

func NewOptions(includePatterns []string) Options {
	return Options{
		IncludePatterns: set.FromSlice(includePatterns),
		ExcludePatterns: set.New[string](),
	}
}

func DefaultGlobalOptions() Options {
	return Options{
		IncludePatterns: set.FromSlice(defaultIncludePatterns),
		ExcludePatterns: set.FromSlice(defaultExcludePatterns),
	}
}

func (o *Options) Merge(other Options) Options {
	return Options{
		IncludePatterns: other.IncludePatterns.Copy().Merge(o.IncludePatterns),
		ExcludePatterns: other.ExcludePatterns.Copy().Merge(o.ExcludePatterns),
	}
}
