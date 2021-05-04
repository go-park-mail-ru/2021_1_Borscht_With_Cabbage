package models

type WsMessageForRepo struct {
	IdMessage  int32  `json:"idMessage"`
	Date       string `json:"date"`
	Content    string `json:"content"`
	SentFromId int32  `json:"sendFromId"`
	SentToId   int32  `json:"sendToId"`
}
