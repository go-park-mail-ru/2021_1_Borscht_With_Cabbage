package models

type WsMessage struct {
	Id   int    `json:"id"`
	Date string `json:"date"`
	Text string `json:"text"`
}

type WsOpponent struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type FromClientPayload struct {
	To      WsOpponent `json:"to"`
	Message WsMessage  `json:"message"`
}

type ToClientPayload struct {
	From    WsOpponent `json:"from"`
	Message WsMessage  `json:"message"`
}

type FromClient struct {
	Action  string            `json:"action"`
	Payload FromClientPayload `json:"payload"`
}

type ToClient struct {
	Action  string          `json:"action"`
	Payload ToClientPayload `json:"payload"`
}

type InfoMessageSend struct {
	From    WsOpponent `json:"from"`
	To      WsOpponent `json:"to"`
	Message WsMessage  `json:"message"`
}

type WsMessageForRepo struct {
	IdMessage  int    `json:"idMessage"`
	Date       string `json:"date"`
	Content    string `json:"content"`
	SentFromId int    `json:"sendFromId"`
	SentToId   int    `json:"sendToId"`
}
