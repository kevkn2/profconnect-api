package port

import "context"

type Service[T, R any] interface {
	Execute(ctx context.Context, input *T) (*R, error)
}
