package repository

import (
	"context"
	"fmt"
	"profconnect-api/internal/database/sqlc/generated"
	"profconnect-api/internal/domain/constants"
	"profconnect-api/internal/domain/entities"
	"profconnect-api/internal/domain/port"

	"github.com/google/uuid"
)

type professorRepository struct {
	queries *generated.Queries
}

func NewProfessorRepository(queries *generated.Queries) port.ProfessorRepository {
	return &professorRepository{
		queries: queries,
	}
}

// Create implements port.ProfessorRepository.
func (p *professorRepository) Create(ctx context.Context, professor *entities.Professor) (*entities.Professor, error) {
	userID, err := uuid.Parse(professor.User.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	newProfessor, err := p.queries.CreateProfessor(
		ctx,
		generated.CreateProfessorParams{
			UserID:     userID,
			University: professor.University,
			Department: professor.Department,
		},
	)
	if err != nil {
		return nil, err
	}

	return &entities.Professor{
		ID:         newProfessor.ID.String(),
		User:       professor.User,
		University: professor.University,
		Department: professor.Department,
	}, nil
}

// GetByID implements port.ProfessorRepository.
func (p *professorRepository) GetByUserID(ctx context.Context, userID string) (*entities.Professor, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	professor, err := p.queries.GetProfessorByUserID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	return &entities.Professor{
		ID: professor.ID.String(),
		User: &entities.User{
			ID:             professor.UserID.String(),
			Name:           professor.UserName,
			Email:          professor.UserEmail,
			HashedPassword: "", // Never expose hashed password
			Role:           constants.Professor,
		},
		University: professor.University,
		Department: professor.Department,
	}, nil
}

// Update implements port.ProfessorRepository.
func (p *professorRepository) Update(ctx context.Context, professor *entities.Professor) (*entities.Professor, error) {
	userID, err := uuid.Parse(professor.User.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	updatedProfessor, err := p.queries.UpdateProfessor(ctx, generated.UpdateProfessorParams{
		UserID:     userID,
		University: professor.University,
		Department: professor.Department,
	})
	if err != nil {
		return nil, err
	}

	return &entities.Professor{
		ID: updatedProfessor.ID.String(),
		User: &entities.User{
			ID:             updatedProfessor.UserID.String(),
			Name:           updatedProfessor.Name,
			Email:          updatedProfessor.Email,
			HashedPassword: "",          // Never expose hashed password
			Role:           "professor", // Assuming role is always professor for this entity
		},
		University: updatedProfessor.University,
		Department: updatedProfessor.Department,
	}, nil
}
