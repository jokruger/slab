package slab

type ClearFunc[T any] func(*T)

// Slab is a simple slab allocator for objects of type T. It maintains an internal buffer of pre-allocated objects
// and allows for efficient allocation and resetting. If the internal buffer is exhausted, it falls back to heap allocation.
type Slab[T any] struct {
	buf      []T
	used     int
	fallback int
	clear    ClearFunc[T]
}

// NewSlab creates a new slab with the given capacity and clear function (a function that releases internal type resources; use nil if not needed).
func NewSlab[T any](size int, clear ClearFunc[T]) Slab[T] {
	return Slab[T]{
		buf:   make([]T, size),
		clear: clear,
	}
}

// Alloc returns a pointer to an object of type T. It first tries to allocate from the slab's internal buffer,
// and if that is exhausted, it falls back to heap allocation.
func (s *Slab[T]) Alloc() *T {
	if s.used < len(s.buf) {
		var zero T
		s.buf[s.used] = zero // zero out the slot before returning it
		t := &s.buf[s.used]
		s.used++
		return t
	}
	s.fallback++
	return new(T) // heap fallback
}

// Reset clears the slab, making all previously allocated objects available for reuse. If a clear function was provided,
// it will be called on each object in the slab before resetting.
func (s *Slab[T]) Reset() {
	if s.clear != nil {
		for i := 0; i < s.used; i++ {
			s.clear(&s.buf[i])
		}
	}
	s.used = 0
	s.fallback = 0
}

// Stats returns the current usage statistics of the slab, including the number of objects allocated from the internal buffer (Pool)
// and the number of objects allocated from the heap (Heap). Stats will be reset to zero after each call to Reset.
func (s *Slab[T]) Stats() Stats {
	return Stats{
		Pool: s.used,
		Heap: s.fallback,
	}
}
