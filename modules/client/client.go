package client

import (
	"Golearn/modules/auth"
	"Golearn/modules/compare"
	"Golearn/modules/container"
	"Golearn/modules/database"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type homePageVars struct {
	Title        string
	Username     string
	IsConnected  bool
	NextExercise database.Exercise
}

type exercicePageVars struct {
	Title          string
	ExerciceTitle  string
	ExercicePrompt string
	ExerciceOutput string
	ExercicesList  []ExerciceLink
	User           string
	IsNotHome      bool
}

type ExerciceLink struct {
	ExerciceName string
	ExerciceDone bool
}

func HomePage(res http.ResponseWriter, req *http.Request) {
	currentLogIn, isOk := auth.ExtractClaims(res, req)
	var nextExercise database.Exercise
	if isOk {
		nextExercise = database.GetExerciseByID(currentLogIn.Progression + 1)
	}
	pageData := homePageVars{
		Title:        "GoLearn | Home",
		Username:     currentLogIn.Name,
		IsConnected:  isOk,
		NextExercise: nextExercise,
	}
	tmpl := template.Must(template.ParseFiles("./CLIENT/static/home.gohtml"))
	err := tmpl.Execute(res, pageData)
	if err != nil {
		log.Println(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func NotLogged(res http.ResponseWriter, req *http.Request) {
	tmpl := template.Must(template.ParseFiles("./CLIENT/static/notLogin.gohtml"))
	_ = tmpl.Execute(res, req)
	return
}

func ExercicePage(res http.ResponseWriter, req *http.Request) {
	var IsNotHome bool
	title := strings.Split(req.URL.Path, "/exercice/")[1]
	currentLogIn, isOk := auth.ExtractClaims(res, req)
	if !isOk {
		http.Redirect(res, req, "/notLogged", http.StatusSeeOther)
		return
	}
	var currentExersise database.Exercise
	if title != "" {
		currentExersise = database.GetExerciseByName(title)
		IsNotHome = true
	} else {
		currentExersise = database.Exercise{IdExercise: 0, Title: "Welcome to GoLearn", Prompt: "Select an exercise to start learning Go !", Difficulty: 0}
		IsNotHome = false
	}
	//PREP Exercises List
	exercicesList := []ExerciceLink{}
	exercices := database.GetExerciseNameList()
	for _, exercice := range exercices {
		if database.GetExerciseByName(exercice).IdExercise <= currentLogIn.Progression {
			exercicesList = append(exercicesList, ExerciceLink{ExerciceName: exercice, ExerciceDone: true})
		} else {
			exercicesList = append(exercicesList, ExerciceLink{ExerciceName: exercice, ExerciceDone: false})
		}
	}
	//END PREP
	pageData := exercicePageVars{
		Title:          "GoLearn | " + currentExersise.Title,
		ExerciceTitle:  currentExersise.Title,
		ExercicePrompt: currentExersise.Prompt,
		ExerciceOutput: "Click Run Code to test you code !",
		ExercicesList:  exercicesList,
		User:           currentLogIn.Name,
		IsNotHome:      IsNotHome,
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
	fmt.Println(jsonResult.Result)
	jsonData, err := json.Marshal(jsonResult)
	if err != nil {
		log.Println(err)
	}
	res.Write(jsonData)
	compare.Compar(compare.GetSolution(req), string(result), res, req)

}
