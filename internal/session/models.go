package session

import "github.com/borscht/backend/internal/models"

type SessionUsecase interface {
	Create(sessionInfo models.SessionInfo) (string, error)
	Check(sessionToCheck string) (models.SessionInfo, bool, error)
	Delete(session string) error
}

type SessionRepo interface {
	Create(sessionData models.SessionData) error
	Check(sessionToCheck string) (models.SessionInfo, bool, error)
	Delete(session string) error
}
