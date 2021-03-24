package domain

type User struct {
	Uid         string `json:"uid"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Phone       string `json:"number"`
	Avatar      string `json:"avatar"`
	MainAddress string `json:"address"`
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
	Create(user User) error
	LoginCheck(user UserAuth) (User, error)
	GetByUid(uid int32) (User, error)
	Update(newUser UserData) error
}

type UserRepo interface {
	Create(user User) error
	CheckUserExists(user UserAuth) (User, error)
	GetByUid(uid int32) (User, error)
	Update(newUser UserData) error
}
