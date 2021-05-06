package models

type Message struct {
	Mid  int32
	Date string
	Text string
}

type User struct {
	Id   int32
	Role string
}

type ChatInfo struct {
	Message Message
	User    User
}

type Chat struct {
	Message   Message
	Sender    User
	Recipient User
}
