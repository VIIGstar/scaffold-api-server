package session_handlers

import (
	"logur.dev/logur"
	"scaffold-api-server/internal/repository"
)

type sessionHandler struct {
	logger logur.LoggerFacade
	repo   repository.Registry
}

func New(logger logur.LoggerFacade, repo repository.Registry) sessionHandler {
	return sessionHandler{
		logger: logger,
		repo:   repo,
	}
}
