package models

type BriefInfoChat struct {
	InfoOpponent
	LastMessage string `json:"last message"`
}

type InfoMessage struct {
	Id     int    `json:"id"`
	Date   string `json:"date"`
	Text   string `json:"text"`
	FromMe bool   `json:"fromMe"`
}

type InfoOpponent struct {
	Uid    int    `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type InfoChat struct {
	InfoOpponent
	Messages []InfoMessage `json:"messages"`
}
