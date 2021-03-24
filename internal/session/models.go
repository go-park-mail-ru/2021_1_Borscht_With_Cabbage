package session

import (
	"context"
)

type SessionUsecase interface {
	Create(ctx context.Context, uid string) (string, error)
	Check(ctx context.Context, sessionToCheck string) (string, bool)
	Update(ctx context.Context, newValue, oldValue string) error // изменяет значение мап по заданной сессии
	Delete(ctx context.Context, session string)
}

type SessionRepo interface {
	Create(ctx context.Context, session, uid string) error
	Check(ctx context.Context, sessionToCheck string) (string, bool)
	Update(ctx context.Context, newValue, oldValue string) error
	Delete(ctx context.Context, session string)
}
