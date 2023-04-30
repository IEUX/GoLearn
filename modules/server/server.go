package server

import (
	"Golearn/modules/auth"
	"Golearn/modules/client"
	"net/http"
)

func InitServer() *http.ServeMux {
	server := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./CLIENT/static/"))
	server.Handle("/static/", http.StripPrefix("/static", fileServer))
	server.HandleFunc("/", client.HomePage)
	server.HandleFunc("/exercice", client.ExercicePage)
	server.HandleFunc("/login", auth.Login)
	server.HandleFunc("/register", auth.Register)
	server.HandleFunc("/logout", auth.Logout)
	return server
}
