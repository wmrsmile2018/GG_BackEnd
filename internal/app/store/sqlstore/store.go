package sqlstore

import (
	"database/sql"
	_ "github.com/lib/pq" // ... анонимный импорт
	"github.com/wmrsmile2018/GG/internal/app/store"
)

type Store struct {
	db             *sql.DB
	userRepository *UserRepository
}

// New
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

//User ...
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}

// store.User().Create()
