package session

type SessionUsecase interface {
	Create(uid int32) (string, error)
	Check(sessionToCheck string) (int32, bool)
	Delete(session string) error
}

type SessionRepo interface {
	Create(session string, uid int32) error
	Check(sessionToCheck string) (int32, bool)
	Delete(session string) error
}
