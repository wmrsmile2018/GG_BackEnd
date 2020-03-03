package store

import (
	"github.com/wmrsmile2018/GG/internal/app/model"
)

//UserRepository ...
type UserRepository interface {
	CreateUser(*model.User) error
	FindByEmail(string) (*model.User, error)
	Find(string) (*model.User, error)
	FindByChat(idChat string) (map[string]bool , error)
	CreateMessage(message *model.Message)(*model.Message, error)
}