package slices

func Contains[T comparable](target T, slice []T) bool {
	for _, sliceVal := range slice {
		if target == sliceVal {
			return true
		}
	}

	return false
}
