package interfaces

type IOption[T any] interface {
	Apply(T)
}

type OptionFunc[T any] func(T)

func (f OptionFunc[T]) Apply(elem T) {
	f(elem)
}
