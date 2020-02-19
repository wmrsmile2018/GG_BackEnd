package model

type Message struct {
	IdMessage		string
	TypeChat		string
	IdChat			string
	User			*User
	IdUser			string
	TimeCreateM		int64
	BytesMessage	[]byte
	Message 		string
}
