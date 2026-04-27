# slab

slab is a small Go package for arena-style allocators built around reusable slabs.

It provides primitives:

- Slab[T] for reusable object slots.
- SliceSlab[T] for reusable slice buffers.

Designed for hot paths where values are created in batches, used briefly, then discarded together, or reused in a cycles:

- Preallocate once.
- Allocate many times from the slab.
- Reset between cycles.
- When limits are reached, allocation falls back to the heap, so behavior stays simple and safe.
