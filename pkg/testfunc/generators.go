package testfunc

func SliceGenerator[T any](quantity uint, creator func() T) []T {
	var res []T
	for i := 0; i < int(quantity); i++ {
		res = append(res, creator())
	}
	return res
}
