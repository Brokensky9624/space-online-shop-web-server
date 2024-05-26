package tool

func FilterSlice[T any](sl []T, filterFunc func(el T) bool) []T {
	newSl := make([]T, 0)
	for _, el := range sl {
		if filterFunc(el) {
			newSl = append(newSl, el)
		}
	}
	return newSl
}
