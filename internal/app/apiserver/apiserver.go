package apiserver

import (
	"database/sql"
	"github.com/gopherschool/http-rest-api/internal/app/store/sqlstore"
	"github.com/gorilla/handlers"
	"github.com/gorilla/sessions"
	"net/http"
	"os"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()
	store := sqlstore.New(db)
	sessionsStore := sessions.NewCookieStore([]byte(config.SessionsKey))
	srv := newServer(store, sessionsStore)
	return http.ListenAndServe(config.BindAddr, handlers.LoggingHandler(os.Stdout, srv))

}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

	
