package collections

func ToAnySlice[T any](in []T) []any {
	ret := make([]any, 0, len(in))

	for _, v := range in {
		ret = append(ret, v)
	}

	return ret
}

func FromAnySlice[T any](in []any) []T {
	ret := make([]T, 0, len(in))

	for _, v := range in {
		vT, ok := v.(T)
		if !ok {
			continue
		}

		ret = append(ret, vT)
	}

	return ret
}
