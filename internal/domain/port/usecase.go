package port

type Usecase[T, R any] interface {
	Execute(input *T) (*R, error)
}
