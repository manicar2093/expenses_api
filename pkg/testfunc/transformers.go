package testfunc

func ToPointer[T any](value T) *T {
	var pointer T = value
	return &pointer
}
