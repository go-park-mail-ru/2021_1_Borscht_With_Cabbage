package models

type User struct {
	Uid         int    `json:"uid"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Phone       string `json:"number"`
	Avatar      string `json:"avatar"`
	MainAddress string `json:"activeAddress"`
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

type SuccessUserResponse struct {
	User
	Role string `json:"role"`
}
