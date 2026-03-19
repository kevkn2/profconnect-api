package port

import (
	"context"
	"profconnect-api/internal/domain/entities"
)

type ProfessorRepository interface {
	GetByUserID(ctx context.Context, id string) (*entities.Professor, error)
	Create(ctx context.Context, professor *entities.Professor) (*entities.Professor, error)
	Update(ctx context.Context, professor *entities.Professor) (*entities.Professor, error)
}
