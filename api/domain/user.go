package domain

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"number"`
	Avatar   string `json:"avatar"`
}

type UserAuth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserUsecase interface {
	Create(ctx *CustomContext, user User) error
	GetByLogin(ctx *CustomContext, user UserAuth) (User, error)
}

type UserRepo interface {
	Create(ctx *CustomContext, user User) error
	GetByLogin(ctx *CustomContext, user UserAuth) (User, error)
}
