package collections

func DefaultIfEmpty[T comparable](v, def T) T {
	var zero T

	if v == zero {
		return def
	}

	return v
}

func IsAnyOf[T comparable](v T, opts ...T) bool {
	for _, opt := range opts {
		if v == opt {
			return true
		}
	}

	return false
}
