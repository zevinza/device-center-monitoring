package lib

// Pointer return pointer of v
// this is more save than directly to &v
func Pointer[T any](v T) *T {
	return &v
}

// Rev return value of pointer v
// if v is nil will return default value in n or default value of its type
func Rev[T any](v *T, n ...T) T {
	var null T
	if len(n) > 0 {
		null = n[0]
	}
	if v == nil {
		return null
	}

	return *v
}

// ComparePtr to compare 2 comparable data in pointer, when one of data is nil, return false
func ComparePtr[T comparable](a, b *T) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// validate that pointer v is not zero value a.k.a "" or 0
func IsValidPtr[T comparable](v *T) bool {
	if v == nil {
		return false
	}

	var zero T
	return *v != zero
}
