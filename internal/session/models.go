package session

type SessionUsecase interface {
	Create(uid int) (string, error)
	Check(sessionToCheck string) (int, bool)
	Delete(session string) error
}

type SessionRepo interface {
	Create(session string, uid int) error
	Check(sessionToCheck string) (int, bool)
	Delete(session string) error
}
