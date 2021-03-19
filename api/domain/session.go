package domain

type Session struct {
	Session string `json:"session"`
	Number  string `json:"number"`
}

type SessionUsecase interface {
	Create(ctx *CustomContext, uid string) (string, error)
	Check(sessionToCheck string, ctx *CustomContext) (string, bool)
}

type SessionRepo interface {
	Create(ctx *CustomContext, session, uid string) error
	Check(sessionToCheck string, ctx *CustomContext) (string, bool)
}
