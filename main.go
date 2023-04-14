package main

import (
	"Golearn/modules/database"
	"Golearn/modules/server"
	"database/sql"
	"fmt"
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
	fmt.Println("Server is running on http://localhost" + port)
	http.ListenAndServe(port, server)
}
