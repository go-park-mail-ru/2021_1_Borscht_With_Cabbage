package session

type SessionUsecase interface {
	Create(uid int, role string) (string, error)
	Check(sessionToCheck string) (int, bool, string)
	Delete(session string) error
}

type SessionRepo interface {
	Create(session string, uid int, role string) error
	Check(sessionToCheck string) (int, bool, string)
	Delete(session string) error
}
