package repository

import (
	"logur.dev/logur"
	"scaffold-api-server/internal/repository/user"
	"scaffold-api-server/internal/services/database"
)

type DatabaseRepo interface {
	User() user.Repo
}

func NewDBImpl(logger logur.LoggerFacade, db *database.DB) dbImpl {
	return dbImpl{
		user: user.New(logger, db),
	}
}

type dbImpl struct {
	user user.Repo
}

func (i dbImpl) User() user.Repo {
	return i.user
}
