package session

import (
	"context"

	"github.com/borscht/backend/internal/models"
)

type SessionUsecase interface {
	Create(ctx context.Context, sessionInfo models.SessionInfo) (string, error)
	Check(ctx context.Context, sessionToCheck string) (models.SessionInfo, bool, error)
	Delete(ctx context.Context, session string) error
	CreateKey(ctx context.Context) (string, error)
	CheckKey(ctx context.Context, sessionToCheck string) (models.SessionInfo, bool, error)
}

type SessionRepo interface {
	Create(ctx context.Context, sessionData models.SessionData) error
	Check(ctx context.Context, sessionToCheck string) (models.SessionInfo, bool, error)
	Delete(ctx context.Context, session string) error
}
