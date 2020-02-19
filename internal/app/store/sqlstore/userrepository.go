package sqlstore

import (
	"database/sql"
	"fmt"
	"github.com/wmrsmile2018/GG/internal/app/model"
	"github.com/wmrsmile2018/GG/internal/app/store"
	"time"
)

//UserRepository ...
type UserRepository struct {
	store *Store
}

//Create...
func (r *UserRepository) CreateUser(u *model.User) (error) {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}
	r.store.db.QueryRow(
		"INSERT INTO users (id_user, email, encrypted_password) VALUES ($1, $2, $3)",
			u.ID,
			u.Email,
			u.EncryptedPassword,
		)
	return nil
}

// CreateMessage
func (r *UserRepository) CreateMessage(message *model.Message) (*model.Message, error) {
	//var id_mes string
	mes := &model.Message{}
	timestamp := time.Unix(0, message.TimeCreateM).Format("2006-01-02, 15:04:05")
	if err := r.store.db.QueryRow(
		"INSERT INTO messages VALUES ($1, $2, $3, $4, $5, $6) RETURNING *",
		message.IdMessage,
		message.IdUser,
		message.IdChat,
		message.Message,
		timestamp,
		message.TypeChat,
		).Scan(
			&mes.IdMessage,
			&mes.IdUser,
			&mes.IdChat,
			&mes.Message,
			&timestamp,
			&mes.TypeChat,
			); err != nil {
				return nil, err
	}
	return mes, nil
}

//FindByMail...
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
 	if err := r.store.db.QueryRow(
		"SELECT id_user, email, encrypted_password FROM users WHERE email = $1",
		email,
		).Scan(
			&u.ID,
			&u.Email,
			&u.EncryptedPassword,
			); err != nil {
				if err == sql.ErrNoRows {
					return nil, store.ErrRecordNotFound
				}
				return nil, err
			}
	return u, nil
}

//Find...
func (r *UserRepository) Find(idUser string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id_user, email, encrypted_password FROM users WHERE id_user = $1",
		idUser,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}


//FincByChat
func (r *UserRepository) FindByChat(idChat string) (map [*model.User]bool, error) {
	mapU := make(map[*model.User]bool)
	u := model.User{}
	rows, err := r.store.db.Query(
		"SELECT id_user, email FROM users WHERE id_user IN (SELECT id_user FROM users_chats WHERE id_chat = $1)", idChat)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(
			&u.ID,
			&u.Email,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}
		fmt.Println("_________________store", u)
		go func(user model.User){mapU[&user] = true}(u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return mapU, nil
}