package client

import (
	"html/template"
	"log"
	"net/http"
)

type exercicePageVars struct {
	Title          string
	ExerciceTitle  string
	ExercicePrompt string
	ExerciceOutput string
	ExercicesList  []ExerciceLink
}

type ExerciceLink struct {
	ID           int
	ExerciceName string
	ExerciceDone bool
}

var exercice1 ExerciceLink = ExerciceLink{
	ID:           1,
	ExerciceName: "Hello World",
	ExerciceDone: true,
}

var exercice2 ExerciceLink = ExerciceLink{
	ID:           2,
	ExerciceName: "Print Alphabet",
	ExerciceDone: false,
}

var exercice3 ExerciceLink = ExerciceLink{
	ID:           3,
	ExerciceName: "Print Numbers",
	ExerciceDone: false,
}

func HomePage(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Home page"))
	res.WriteHeader(http.StatusOK)
}

func ExercicePage(res http.ResponseWriter, req *http.Request) {
	var exercicesList []ExerciceLink
	exercicesList = append(exercicesList, exercice1, exercice2, exercice3)
	pageData := exercicePageVars{
		Title:          "GoLearn | Hello World !",
		ExerciceTitle:  "Hello World ! 🌎",
		ExercicePrompt: "Write a program that prints ‘Hello World’ to the screen.",
		ExerciceOutput: "Click <a>Run Code</a> to test you code ",
		ExercicesList:  exercicesList,
	}
	tmpl := template.Must(template.ParseFiles("./CLIENT/static/exercicePage.gohtml"))
	err := tmpl.Execute(res, pageData)
	if err != nil {
		log.Println(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func CORSManager(res http.ResponseWriter, req *http.Request) {
	origin := req.Header.Get("Origin")
	res.Header().Set("Access-Control-Allow-Origin", origin)
	res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	res.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	res.Header().Set("Access-Control-Allow-Credentials", "true")
}
