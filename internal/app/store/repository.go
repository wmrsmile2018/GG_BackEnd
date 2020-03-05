package store

import (
	"github.com/wmrsmile2018/GG/internal/app/model"
)

//UserRepository ...
type UserRepository interface {
	CreateUser(*model.User) (*model.User, error)
	FindByEmail(id string) (*model.User, error)
	Find(id string) (*model.User, error)
	FindByChat(idChat string) (map[string]bool, error)
	CreateMessage(message *model.Message) (*model.Message, error)
	PaginationMessages(params *model.ParametersPagination) ([]model.Message, error)
	CreateChat(idChat string, idUser string, typeChat string) (*model.Chat, error)
	CreateUserChat(idChat string, idUser string) (*model.UserChat, error)
}
