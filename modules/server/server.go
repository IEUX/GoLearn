package server

import (
	"Golearn/modules/application"
	"Golearn/modules/auth"
	"Golearn/modules/client"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func InitServer() *mux.Router {
	server := mux.NewRouter()
	server.HandleFunc("/", client.HomePage)
	server.HandleFunc("/login", auth.Login)
	server.HandleFunc("/register", auth.Register)
	server.HandleFunc("/logout", auth.Logout)
	server.HandleFunc("/recivecode", application.ReciveCode)
	server.HandleFunc("/helloworld", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello world")
	})
	return server
}
