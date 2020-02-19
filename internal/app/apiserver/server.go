package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/wmrsmile2018/GG/internal/app/model"
	"github.com/wmrsmile2018/GG/internal/app/service"
	"github.com/wmrsmile2018/GG/internal/app/store"
	"net/http"
	"time"
)

var (
	errIncorrecntingEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated = errors.New("not authenticated")
)


var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	sessionName = "hellomydear"
	ctxKeyUser ctxKey = iota
	ctxKeyRequestID
)

type ctxKey int8

type server struct {
	router			*mux.Router
	logger			*logrus.Logger
	store			store.Store
	sessionsStore	sessions.Store
	hub				*service.Hub
}

func newServer(hub *service.Hub, store store.Store, sessionsStore sessions.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
		sessionsStore: sessionsStore,
		hub: hub,
	}
	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server)  configureRouter(){
	s.router.Use(s.setRequestID) // подставляет для каждого входящего запроса уникальный id для передачи в заголовки и для лога
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string {"*"}))) // если запрос приходит с другого порта
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/sessions",s.handleSessionsCreate()).Methods("POST")
	s.router.HandleFunc("/general_chat", s.hadnleUsersGeneralChat()).Methods("GET")
	s.router.HandleFunc("/local_chat", s.hadnleUsersLocalChat()).Methods("GET")
	s.router.HandleFunc("/ws", s.handleWS())
	s.router.HandleFunc("/test", s.handleTest()).Methods("GET")
	//private/***
	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/whoami", s.handleWhoAmI()).Methods("GET")
}

func (s *server) handleTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := &model.Message{
			IdMessage:    "a0eebc19-9c0b-4ef8-bb6d-6bb9bd380a10",
			TypeChat:     "twosome",
			IdChat:       "a0eebc79-9c0b-4ef8-bb6d-6bb9bd380a10",
			User:         nil,
			IdUser:       "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a10",
			TimeCreateM:  0,
			BytesMessage: nil,
			Message:      "pewpew",
		}
		mes, err := s.store.User().CreateMessage(message)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}
		s.respond(w, r, http.StatusCreated, message)
		fmt.Println(mes)

	}
}

func (s *server) handleWS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		err = r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().Find(r.Form.Get("id"))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		client:= service.NewClient(s.hub, conn, make(chan *model.Message, 256), u)
		clientConn := &service.ClientConn{
			Connection: client,
			Id:			u.ID,
		}
		s.hub.Register <- clientConn
		s.hub.UserConnection[u.ID] = clientConn
		go clientConn.WritePump()
		go clientConn.ReadPump()
	}
}

func (s *server) hadnleUsersGeneralChat() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "home.html")
	}
}

func (s *server) hadnleUsersLocalChat() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "home1.html")
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		params := r.Form
		s.hub.Users, err = s.store.User().FindByChat(params.Get("id"))
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
	}
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id": r.Context().Value(ctxKeyRequestID),
		})

		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWrite{w, http.StatusOK}
		next.ServeHTTP(w, r)

		logger.Infof(
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
			)
	})
}


func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionsStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		u, err := s.store.User().Find(id.(string))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))

	})
}

func (s *server) handleWhoAmI() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*model.User))
	}
}


func (s *server) handleSessionsCreate() http.HandlerFunc {
	type request struct { // ожидаемые поля от пользователя при регистрации
		Email		string 	`json:"email"`
		Password	string 	`json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u, err := s.store.User().FindByEmail(req.Email)

		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrecntingEmailOrPassword)
			return
		}

		session, err := s.sessionsStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		session.Values["user_id"] = u.ID
		if err := s.sessionsStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleUsersCreate() http.HandlerFunc {
	type request struct { // ожидаемые поля от пользователя при регистрации
		Email		string 	`json:"email"`
		Password	string 	`json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u := &model.User{
			Email:				req.Email,
			Password:			req.Password,
			ID:					uuid.New().String(),
		}
		if err := s.store.User().CreateUser(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
		}
		u.Sanitize()
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error__": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			return
		}

	}
}

