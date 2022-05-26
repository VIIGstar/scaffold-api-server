package user

import (
	"context"
	"logur.dev/logur"
	"scaffold-api-server/internal/entities"
	query_params "scaffold-api-server/internal/query-params"
	"scaffold-api-server/internal/services/database"
)

// New creates new impl impl and returns as User interface
func New(logger logur.LoggerFacade, db *database.DB) Repo {
	return &impl{
		logger: logger,
		db:     db,
	}
}

// Repo represents methods that User repository must implement
type Repo interface {
	// Create inserts new record in User table
	Create(u *entities.Investor) error
	// Find retrieves a impl based on search criteria
	Find(ctx context.Context, req query_params.GetUserParams, lock bool) (*entities.Investor, error)
}

type impl struct {
	logger logur.LoggerFacade
	db     *database.DB
}
