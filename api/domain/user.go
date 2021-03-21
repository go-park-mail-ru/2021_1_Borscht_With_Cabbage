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

type UserData struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"number"`
	Password    string `json:"password"`
	PasswordOld string `json:"password_current"`
	Avatar      string `json:"avatar"`
}

type UserReg struct {
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type SuccessResponse struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type UserUsecase interface {
	Create(ctx *CustomContext, user User) error
	GetByLogin(ctx *CustomContext, user UserAuth) (User, error)
	GetByNumber(ctx *CustomContext, number string) (User, error)
	Update(ctx *CustomContext, newUser UserData) error
}

type UserRepo interface {
	Create(ctx *CustomContext, user User) error
	GetByLogin(ctx *CustomContext, user UserAuth) (User, error)
	GetByNumber(ctx *CustomContext, number string) (User, error)
	Update(ctx *CustomContext, newUser UserData) error
}
