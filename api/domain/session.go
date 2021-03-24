package domain

type Session struct {
	Session string `json:"session"`
	Number  string `json:"number"`
}

type SessionUsecase interface {
	Create(uid string) (string, error)
	Check(sessionToCheck string) (int32, bool)
	UpdateValue(newValue, oldValue string) error // изменяет значение мап по заданной сессии
	Delete(session string) error
}

type SessionRepo interface {
	Create(session, uid string) error
	Check(sessionToCheck string) (int32, bool)
	UpdateValue(newValue, oldValue string) error
	Delete(session string) error
}
