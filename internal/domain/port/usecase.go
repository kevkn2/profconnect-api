package port

import "context"

type Usecase[T, R any] interface {
	Execute(ctx context.Context, input *T) (*R, error)
}
