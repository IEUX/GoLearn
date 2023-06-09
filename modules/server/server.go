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
	server.HandleFunc("/exercice/", client.ExercicePage)
	server.HandleFunc("/login", auth.Login)
	server.HandleFunc("/signup", auth.Register)
	server.HandleFunc("/logout", auth.Logout)
	server.HandleFunc("/sendCode", client.SendCode)
	server.HandleFunc("/401", client.NotLogged)
	server.HandleFunc("/favicon.ico", client.Favicon)
	return server
}
