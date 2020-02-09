package main

import (
  "github.com/gorilla/handlers"
  "github.com/gorilla/mux"
  "net/http"
  "os"
)

type server struct {
  router *mux.Router

}


func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  s.router.ServeHTTP(w, r)
}

func newServer () *server {
  s := &server{
    router: mux.NewRouter(),
  }

  s.configureRouter()

  return s
}

func (s *server) configureRouter (){
  s.router.HandleFunc("/", s.handleHello())
}

func (s *server) handleHello() http.HandlerFunc {
  return func( w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("hello"))
  }
}

func main() {
  srv := newServer()
  http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, srv))
}
