package set

type Set[T comparable] map[T]bool

func New[T comparable]() Set[T] {
	return make(Set[T])
}

func FromSlice[T comparable](slice []T) Set[T] {
	s := make(Set[T])
	for _, value := range slice {
		s[value] = true
	}
	return s
}

func (s Set[T]) Add(value ...T) {
	for _, v := range value {
		s[v] = true
	}
}

func (s Set[T]) Contains(value T) bool {
	_, ok := s[value]
	return ok
}

func (s Set[T]) Remove(value T) {
	delete(s, value)
}

func (s Set[T]) Copy() Set[T] {
	c := make(Set[T])
	for value := range s {
		c[value] = true
	}
	return c
}

func (s Set[T]) Merge(other Set[T]) Set[T] {
	for value := range other {
		s[value] = true
	}
	return s
}

func (s Set[T]) Slice() []T {
	var slice []T
	for value := range s {
		slice = append(slice, value)
	}
	return slice
}
