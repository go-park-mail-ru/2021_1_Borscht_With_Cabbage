package session

import "context"

type SessionUsecase interface {
	Create(ctx context.Context, uid int, role string) (string, error)
	Check(ctx context.Context, sessionToCheck string) (int, bool, string)
	Delete(ctx context.Context, session string) error
}

type SessionRepo interface {
	Create(ctx context.Context, session string, uid int, role string) error
	Check(ctx context.Context, sessionToCheck string) (int, bool, string)
	Delete(ctx context.Context, session string) error
}
