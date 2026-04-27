package slab

// SliceSlab is a simple slab allocator for slices. It preallocates a number of slices with a fixed capacity and returns them on demand.
// If the requested capacity exceeds the preallocated capacity or if all preallocated slices are used, it falls back to allocating a new slice on the heap.
type SliceSlab[T any] struct {
	buf      [][]T
	cap      int
	used     int
	fallback int
}

// NewSliceSlab creates a new SliceSlab with the given number of slices (size) and capacity for each slice (cap).
func NewSliceSlab[T any](size int, cap int) SliceSlab[T] {
	buf := make([][]T, size)
	for i := range buf {
		buf[i] = make([]T, 0, cap)
	}
	return SliceSlab[T]{buf: buf, cap: cap}
}

// Alloc returns a slice of type T with the requested minimum capacity. It first tries to allocate from the slab's internal buffer,
// and if that is exhausted or if the requested capacity exceeds the preallocated capacity, it falls back to heap allocation.
// If sized is true, the returned slice will have a length equal to the requested capacity; otherwise, it will have a length of 0 and the requested capacity.
func (s *SliceSlab[T]) Alloc(cap int, sized bool) []T {
	if cap <= s.cap && s.used < len(s.buf) {
		t := s.buf[s.used]
		if sized {
			t = t[:cap]
		}
		s.used++
		return t
	}
	s.fallback++
	if sized {
		return make([]T, cap)
	}
	return make([]T, 0, cap)
}

// Reset clears the slab, making all previously allocated slices available for reuse.
func (s *SliceSlab[T]) Reset() {
	// No need to clear the slices since they never modified (we modify the copy in alloc, not the original) and cannot exceed the original capacity.
	s.used = 0
	s.fallback = 0
}
