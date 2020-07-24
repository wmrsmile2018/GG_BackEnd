package sqlstore

import (
	"database/sql"
	"github.com/wmrsmile2018/GG/internal/app/model"
	"github.com/wmrsmile2018/GG/internal/app/store"
	"time"
)

//UserRepository ...
type UserRepository struct {
	store *Store
}

//CreateUser - create new user
func (r *UserRepository) CreateUser(u *model.User) (*model.User, error) {
	var sU model.User
	if err := u.Validate(); err != nil {
		return nil, err
	}
	if err := u.BeforeCreate(); err != nil {
		return nil, err
	}
	if err := r.store.db.QueryRow(
		"INSERT INTO users (id_user, email, encrypted_password) VALUES ($1, $2, $3) RETURNING id_user, email, encrypted_password",
		u.ID,
		u.Email,
		u.EncryptedPassword,
	).Scan(
		&sU.ID,
		&sU.Email,
		&sU.EncryptedPassword,
		); err != nil {
			return nil, err
	}
	return &sU, nil
}

//CreateChat - create new unique chat
func (r *UserRepository) CreateChat(idChat string, idUser string, typeChat string) (*model.Chat, error) {
	var c model.Chat
	if err := r.store.db.QueryRow(
		"INSERT INTO chats (id_chat, id_user, type_chat) VALUES ($1, $2, $3) RETURNING id_chat, id_user, type_chat",
		idChat,
		idUser,
		typeChat,
	).Scan(
		&c.IdChat,
		&c.IdUser,
		&c.TypeChat,
		); err != nil {
		return nil, err
	}
	return &c, nil
}

//CreateUserChat - create new records about users inside 1 chat
func (r *UserRepository) CreateUserChat(idChat string, idUser string) (*model.UserChat, error) {
	var c model.UserChat
	if err := r.store.db.QueryRow(
		"INSERT INTO users_chats (id_chat, id_user) VALUES ($1, $2) RETURNING id_chat, id_user",
		idChat,
		idUser,
	).Scan(
		&c.IdChat,
		&c.IdUser,
	); err != nil {
		return nil, err
	}
	return &c, nil
}


// CreateMessage
func (r *UserRepository) CreateMessage(message *model.Message) (*model.Message, error) {
	mes := &model.Message{}
	var number int
	timeStamp := time.Unix(message.TimeCreateM, 0).Format(time.RFC3339)
	if err := r.store.db.QueryRow(
		"INSERT INTO messages (id_message, id_user, id_chat, text_mes, creation_time, type_chat) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *",
		message.IdMessage,
		message.IdUser,
		message.IdChat,
		message.Message,
		timeStamp,
		message.TypeChat,
	).Scan(
		&number,
		&mes.IdMessage,
		&mes.IdUser,
		&mes.IdChat,
		&mes.Message,
		&timeStamp,
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

//func (r *UserRepository) PaginationMessages(params *model.ParametersPagination) (*model.Message, error) {
//	m := &model.Message{}
//	var timeStamp string
//	if err := r.store.db.QueryRow(
//		"SELECT id_user, text_mes, creation_time FROM messages WHERE NUMBER > $1 AND id_chat = $2 limit $3",
//		params.Where,
//		params.Id,
//		params.Number,
//	).Scan(
//		&m.IdUser,
//		&m.Message,
//		&timeStamp,
//	); err != nil {
//		if err == sql.ErrNoRows {
//			return nil, store.ErrRecordNotFound
//		}
//		return nil, err
//	}
//
//	t, err := time.Parse(time.RFC3339, timeStamp)
//	if err != nil {
//		return nil, err
//	}
//	m.TimeCreateM = t.Unix()
//	return m, nil
//}

//PaginationMessages...
func (r *UserRepository) PaginationMessages(params *model.ParametersPagination) ([]model.Message, error) {
	var sMes []model.Message
	rows, err := r.store.db.Query(
		"SELECT id_user, text_mes, creation_time FROM messages WHERE NUMBER > $1 AND id_chat = $2 limit $3",
		params.Where,
		params.Id,
		params.Number,
	)
	if err != nil {

	}
	defer rows.Close()
	for rows.Next() {
		m := model.Message{}
		var timeStamp string
		if err = rows.Scan(
			&m.IdUser,
			&m.Message,
			&timeStamp,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}
		t, err := time.Parse(time.RFC3339, timeStamp)
		if err != nil {
			return nil, err
		}
		m.TimeCreateM = t.Unix()
		sMes = append(sMes, m)
	}
	return sMes, nil
}

//FindByChat
func (r *UserRepository) FindByChat(idChat string) (map[string]bool, error) {
	mapU := make(map[string]bool)
	rows, err := r.store.db.Query(
		"SELECT id_user FROM users WHERE id_user IN (SELECT id_user FROM users_chats WHERE id_chat = $1)",
		idChat,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var UId string
		if err = rows.Scan(
			&UId,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, store.ErrRecordNotFound
			}
			return nil, err
		}
		mapU[UId] = true
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return mapU, nil
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
