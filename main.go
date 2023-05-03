package main

import (
	"Golearn/modules/database"
	"Golearn/modules/server"
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	Router   *mux.Router
	Database *sql.DB
}

func main() {
	cnx := database.GetDbInstance()
	defer cnx.Close()
	server := server.InitServer()
	port := ":8080"
	log.Println("[SERVER] Server is running on http://localhost" + port)
	err := http.ListenAndServe(port, server)
	if err != nil {
		log.Fatalln(err)
	}
}
