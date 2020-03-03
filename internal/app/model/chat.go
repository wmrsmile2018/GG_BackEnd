package model

type Message struct {
	IdMessage		string
	TypeChat		string
	IdChat			string
	IdUser			string
	TimeCreateM		int64
	BytesMessage	[]byte
	Message 		string
}

type Send struct {
	User	*User
	Message *Message
}
