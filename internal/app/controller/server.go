package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/gddo/httputil/header"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/wmrsmile2018/GG/internal/app/model"
	"github.com/wmrsmile2018/GG/internal/app/service"
	"github.com/wmrsmile2018/GG/internal/app/store"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	errIncorrecntingEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated             = errors.New("not authenticated")
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	sessionName        = "hellomydear"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID
)

type ctxKey int8

type server struct {
	router        *mux.Router
	logger        *logrus.Logger
	store         store.Store
	sessionsStore sessions.Store
	hub           *service.Hub
}

type malformedRequest struct {
	status int
	msg    string
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}

func newServer(hub *service.Hub, store store.Store, sessionsStore sessions.Store) *server {
	s := &server{
		router:        mux.NewRouter(),
		logger:        logrus.New(),
		store:         store,
		sessionsStore: sessionsStore,
		hub:           hub,
	}
	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID) // подставляет для каждого входящего запроса уникальный id для передачи в заголовки и для лога
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"}))) // если запрос приходит с другого порта
	s.router.HandleFunc("/createChat", s.handleCreateChat()).Methods("POST")
	s.router.HandleFunc("/createUserChat", s.handleCreateUserChat()).Methods("POST")
	s.router.HandleFunc("/createChat", s.handleCreateUserChat()).Methods("POST")
	s.router.HandleFunc("/getMessages", s.handleGetMessage()).Methods("GET")
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST")
	s.router.HandleFunc("/ws", s.handleWS())
	s.router.HandleFunc("/test", s.handleTest()).Methods("GET")

	//chats/***
	chat := s.router.PathPrefix("/chat").Subrouter()
	chat.HandleFunc("/general", s.handleUsersGeneralChat()).Methods("GET")
	chat.HandleFunc("/local", s.handleUsersLocalChat()).Methods("GET")
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
			IdChat:       "a0eebc79-9c0b-4ef8-bb6d-6bb9bd380a14",
			IdUser:       "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a10",
			TimeCreateM:  0,
			BytesMessage: nil,
			Message:      "pewpew",
		}
		fmt.Println(s.store.User().CreateChat(message.IdChat, message.IdUser, message.TypeChat))
		fmt.Println(s.store.User().CreateUserChat(message.IdChat, message.IdUser))
	}
}

func (s *server) handleWS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if err = r.ParseForm(); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().Find(r.Form.Get("id"))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		client := service.NewClient(s.hub, ws, u, make(chan *model.Send, 512))
		s.hub.Register <- client
		if len(s.hub.AllConns[u.ID]) == 0 {
			s.hub.AllConns[u.ID] = make(map[*service.ChatClients]bool)
		}
		conn := s.hub.AllConns[u.ID]
		conn[client] = true

		go client.WritePump()
		go client.ReadPump()
	}
}

//curl -v -H "Content-Type: application/json" -X POST \
//-d '{"idChat":"a0eebc79-9c0b-4ef8-bb6d-6bb9bd380a16", "idUser":"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a10", "TypeChat": "twosome"}' http://localhost:8000/createChat
func (s *server) handleCreateChat() http.HandlerFunc {
	type Chat struct {
		IdChat string
		IdUser string
		TypeChat string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var chat Chat
		err := decodeJSONBody(w, r, &chat)
		if err != nil {
			var mr *malformedRequest
			s.error(w, r, mr.status, err)
			if errors.As(err, &mr) {
				http.Error(w, mr.msg, mr.status)
			} else {
				log.Println(err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}
		res, err := s.store.User().CreateChat(chat.IdChat, chat.IdUser, chat.TypeChat)
		fmt.Println(res, err)
		fmt.Fprintf(w, " \n chat %+v  %+v", res, err)

	}
}

//curl -v -H "Content-Type: application/json" -X POST \
//-d '{"idChat":"your name","idUsers":["hello", "pewpew"]}' http://localhost:8000/createChat
func (s *server) handleCreateUserChat() http.HandlerFunc {
	type UsersChats struct {
		IdChat  string
		IdUsers []string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var usersChats UsersChats
		err := decodeJSONBody(w, r, &usersChats)
		if err != nil {
			var mr *malformedRequest
			s.error(w, r, mr.status, err)
			if errors.As(err, &mr) {
				http.Error(w, mr.msg, mr.status)
			} else {
				log.Println(err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}
		fmt.Fprintf(w, "UsersChats:  %+v", usersChats)
		for idUser := range usersChats.IdUsers {
			res, err := s.store.User().CreateUserChat(usersChats.IdChat, usersChats.IdUsers[idUser])
			fmt.Println(res, err)
			fmt.Fprintf(w, "\n %+v  %+v", res, err)
		}
	}
}

func (s *server) handleUsersGeneralChat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "home.html")
	}
}

func (s *server) handleUsersLocalChat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "home1.html")
	}
}

func (s *server) handleGetMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		body := r.Form
		var params = model.ParametersPagination{
			Id:     body.Get("id"),
			Number: s.atoi(w, r, http.StatusBadRequest, body.Get("number")),
			Where:  s.atoi(w, r, http.StatusBadRequest, body.Get("where")),
		}
		sMes, err := s.store.User().PaginationMessages(&params)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
		}
		data, err := json.Marshal(sMes)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}
		w.Write(data)
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
			"request_id":  r.Context().Value(ctxKeyRequestID),
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

func (s *server) handleWhoAmI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*model.User))
	}
}

func (s *server) handleSessionsCreate() http.HandlerFunc {
	type request struct { // ожидаемые поля от пользователя при регистрации
		Email    string `json:"email"`
		Password string `json:"password"`
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
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
			ID:       uuid.New().String(),
		}
		if _, err := s.store.User().CreateUser(u); err != nil {
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

func (s *server) atoi(w http.ResponseWriter, r *http.Request, code int, value string) int {
	res, err := strconv.Atoi(value)
	if err != nil {
		s.error(w, r, code, err)
		return 0
	}
	return res
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

		default:
			return err
		}
	}

	if dec.More() {
		msg := "Request body must only contain a single JSON object"
		return &malformedRequest{status: http.StatusBadRequest, msg: msg}
	}

	return nil
}
