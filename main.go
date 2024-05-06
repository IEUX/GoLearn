package main

import (
	"Golearn/modules/database"
	"Golearn/modules/server"
	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type App struct {
	Router   *mux.Router
	Database *sql.DB
}

func main() {
	cnx := database.GetDbInstance()
	defer func(cnx *sql.DB) {
		err := cnx.Close()
		if err != nil {

		}
	}(cnx)
	router := server.InitServer()
	port := ":8080"
	log.Println("[SERVER] Server is running on http://localhost" + port)
	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Fatalln(err)
	}
}
