package models

type Session struct {
	Session string `json:"session"`
	Number  string `json:"number"`
}

type Auth struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Role   string `json:"role"`
}

type SessionInfo struct {
	Id   int    `json:"uid"`
	Role string `json:"role"`
}

type SessionData struct {
	Session string `json:"session"`
	Id      int    `json:"uid"`
	Role    string `json:"role"`
}
