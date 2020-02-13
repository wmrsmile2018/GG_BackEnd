package sqlstore

import (
	"database/sql"
	"github.com/wmrsmile2018/GG/internal/app/model"
	"github.com/wmrsmile2018/GG/internal/app/store"
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
	return r.store.db.QueryRow(
		"INSERT INTO users (id_user, email, encrypted_password) VALUES ($1, $2, $3) RETURNING id_user",
			u.ID,
			u.Email,
			u.EncryptedPassword,
		).Scan(&u.ID)

}

////
//func (r *UserRepository) CreateUser(u *model.User) (error) {
//	if err := u.Validate(); err != nil {
//		return err
//	}
//
//	if err := u.BeforeCreate(); err != nil {
//		return err
//	}
//	return r.store.db.QueryRow(
//		"INSERT INTO users (id_user, email, encrypted_password) VALUES ($1, $2, $3) RETURNING id_user",
//		u.ID,
//		u.Email,
//		u.EncryptedPassword,
//	).Scan(&u.ID)
//
//}

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
func (r *UserRepository) Find(id_user string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id_user, email, encrypted_password FROM users WHERE id_user = $1",
		id_user,
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
func (r *UserRepository) FindByChat(idChat string) (chan *model.User, error) {
	chU := make(chan *model.User)
	u := model.User{}
	rows, err := r.store.db.Query(
		"SELECT id_user, email FROM users WHERE id_user IN (SELECT id_user FROM chats WHERE id_chat = $1)", idChat)
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

		go func(u model.User){chU <- &u}(u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return chU, nil
}