package apiserver

import (
	"database/sql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/sessions"
	"github.com/wmrsmile2018/GG/internal/app/service"
	"github.com/wmrsmile2018/GG/internal/app/store/sqlstore"
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
	hub := service.NewHub(store)
	srv := newServer(hub, store, sessionsStore)
	go hub.Run()

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

	
