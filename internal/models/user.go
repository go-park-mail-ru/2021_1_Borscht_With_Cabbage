package models

type User struct {
	Uid          int    `json:"uid"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	HashPassword []byte
	Phone        string  `json:"number"`
	Avatar       string  `json:"avatar"`
	Address      Address `json:"address"`
}

type UserAuth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserImageResponse struct {
	Filename string `json:"avatar"`
}

type UserData struct {
	ID          int     `json"id"`
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	Phone       string  `json:"number"`
	Password    string  `json:"password"`
	PasswordOld string  `json:"password_current"`
	Avatar      string  `json:"avatar"`
	Address     Address `json:"address"`
}

type SuccessUserResponse struct {
	User
	Role string `json:"role"`
}
