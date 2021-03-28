package session

type SessionUsecase interface {
	Create(uid int32, role string) (string, error)
	Check(sessionToCheck string) (int32, bool, string)
	Delete(session string) error
}

type SessionRepo interface {
	Create(session string, uid int32, role string) error
	Check(sessionToCheck string) (int32, bool, string)
	Delete(session string) error
}
