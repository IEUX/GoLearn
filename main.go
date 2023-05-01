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

	// start := time.Now()
	// code := "package main\nimport \"fmt\"\nfunc main() {\nfmt.Println(\"Hello World!\")\n}"
	// userFolder := container.CreateCodeFile("USER_", code)
	// userResponse := container.TestCode(userFolder)
	// err := os.RemoveAll(userFolder)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println("test done in " + time.Since(start).String())

	cnx := database.GetDbInstance()
	defer cnx.Close()
	server := server.InitServer()
	port := ":8080"
	fmt.Println("Server is running on http://localhost" + port)
	err := http.ListenAndServe(port, server)
	if err != nil {
		log.Fatalln(err)
	}
}
