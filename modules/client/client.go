package client

import (
	"Golearn/modules/container"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type exercicePageVars struct {
	Title          string
	ExerciceTitle  string
	ExercicePrompt string
	ExerciceOutput string
	ExercicesList  []ExerciceLink
	User           string
	IsConnected    bool
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
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Home page"))
}

func ExercicePage(res http.ResponseWriter, req *http.Request) {
	var isConnected bool
	c, err := req.Cookie("token")
	if err != nil || c == nil {
		isConnected = false
	} else {
		isConnected = true
	}
	title := strings.Split(req.URL.Path, "/exercice/")[1]
	var exercicesList []ExerciceLink
	exercicesList = append(exercicesList, exercice1, exercice2, exercice3)
	pageData := exercicePageVars{
		Title:          "GoLearn | Hello World !",
		ExerciceTitle:  title,
		ExercicePrompt: "Write a program that prints ‘Hello World’ to the screen.",
		ExerciceOutput: "Click <a>Run Code</a> to test you code ",
		ExercicesList:  exercicesList,
		User:           "User",
		IsConnected:    isConnected,
	}
	tmpl := template.Must(template.ParseFiles("./CLIENT/static/exercicePage.gohtml"))
	err = tmpl.Execute(res, pageData)
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

type Code struct {
	Code string
}

type Result struct {
	Result string
}

func SendCode(res http.ResponseWriter, req *http.Request) {
	var code Code
	b, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(b, &code)
	if err != nil {
		log.Fatalln(err)
	}
	path := container.CreateCodeFile("user", code.Code)
	result := container.TestCode(path)
	os.RemoveAll(path)
	jsonResult := Result{
		Result: string(result),
	}
	fmt.Println((jsonResult.Result))
	jsonData, err := json.Marshal(jsonResult)
	if err != nil {
		log.Println(err)
	}
	res.Write(jsonData)
}
