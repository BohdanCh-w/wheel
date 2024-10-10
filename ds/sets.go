package ds

type Set[T any] interface {
	Add(values ...T)
	Del(values ...T)
	Empty() bool
	Has(val T) bool
	HasAny(values ...T) bool
	HasAll(values ...T) bool
	Len() int
	Each(fn func(T))
	Values() []T
	Clear()
	// next functions use concrete type as argument or return type
	// Map(fn func(T) T) Set[T]
	// Copy() Set[T]
	// Equal(other Set[T]) bool
	// Union(other Set[T]) Set[T]
	// Intersect(other Set[T]) Set[T]
	// Diff(other Set[T]) Set[T]
	// SymDiff(other Set[T]) Set[T]
	// IsSubset(other Set[T]) Set[T]
}
