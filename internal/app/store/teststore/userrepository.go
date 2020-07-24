package teststore

import (
	"github.com/wmrsmile2018/GG/internal/app/model"
	"github.com/wmrsmile2018/GG/internal/app/store"
	"strconv"
)

//UserRepository ...
type UserRepository struct {
	store *Store
	users map[string]*model.User
}

func (r * UserRepository) CreateUser(u *model.User) (*model.User, error) {
	if err := u.Validate(); err != nil {
		return nil, err
	}
	if err := u.BeforeCreate(); err != nil {
		return nil, err
	}
	u.ID = strconv.Itoa(len(r.users) + 1)
	r.users[u.ID] = u
	return u, nil
}

func (r *UserRepository) Find(idUser string) (*model.User, error) {
	u, ok := r.users[idUser]
	if !ok {
		return nil, store.ErrRecordNotFound
	}
	return u, nil
}

func (r *UserRepository) FindByChat(idChat string) (map[string]bool, error) {
	panic("implement me")
}

func (r *UserRepository) CreateMessage(message *model.Message) (*model.Message, error) {
	panic("implement me")
}

func (r *UserRepository) PaginationMessages(params *model.ParametersPagination) ([]model.Message, error){
	panic("implement me")
}

func (r *UserRepository) CreateChat(idChat string, idUser string, typeChat string) (*model.Chat, error) {
	panic("implement me")
}

func (r *UserRepository) CreateUserChat(idChat string, idUser string) (*model.UserChat, error) {
	panic("implement me")
}

//FindByEmail ...
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}

	return nil, store.ErrRecordNotFound
}