package models

type BriefInfoChat struct {
	Uid         int    `json:"id"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	LastMessage string `json:"last message"`
}
