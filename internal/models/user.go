package models

type User struct {
	Uid         int    `json:"uid"`
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

type SuccessResponse struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}
