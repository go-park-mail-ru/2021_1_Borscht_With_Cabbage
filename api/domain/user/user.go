package user

import "backend/api/domain"

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
	Create(ctx *domain.CustomContext, user User) error
	GetByLogin(ctx *domain.CustomContext, user UserAuth) (User, error)
	GetByNumber(ctx *domain.CustomContext, number string) (User, error)
	Update(ctx *domain.CustomContext, newUser UserData) error
}

type UserRepo interface {
	Create(ctx *domain.CustomContext, user User) error
	GetByLogin(ctx *domain.CustomContext, user UserAuth) (User, error)
	GetByNumber(ctx *domain.CustomContext, number string) (User, error)
	Update(ctx *domain.CustomContext, newUser UserData) error
}
