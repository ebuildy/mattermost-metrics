package utils

func SafeDeref[T any](p *T) T {
	if p == nil {
		var v T
		return v
	}
	return *p
}

func SafeRef[T any](p T) *T {
	return &p
}
