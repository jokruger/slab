package slab

type Stats struct {
	Pool int // number of objects used from pool
	Heap int // number of objects allocated on heap
}
