package main

import (
	"Golearn/modules/database"
	"Golearn/modules/server"
	"database/sql"
	"fmt"
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
	port := ":8090"
	fmt.Println("Server is running on http://localhost" + port)
	err := http.ListenAndServe(port, server)
	if err != nil {
		log.Fatal(err)
	}
}
