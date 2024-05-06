package client

import (
	"Golearn/modules/auth"
	"Golearn/modules/compare"
	"Golearn/modules/container"
	"Golearn/modules/database"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type homePageVars struct {
	Title        string
	Username     string
	IsConnected  bool
	NextExercise database.Exercise
}

type exercicePageVars struct {
	Title             string
	ExerciceTitle     string
	ExercicePrompt    string
	ExerciceOutput    string
	FunctionStructure string
	ExercicesList     []ExerciceLink
	CanDo             bool
	User              string
	IsNotHome         bool
}

type ExerciceLink struct {
	ExerciceName string
	ExerciceDone bool
}

func HomePage(res http.ResponseWriter, req *http.Request) {
	currentLogIn, isOk := auth.ExtractClaims(res, req)
	if req.URL.Path != "/" {
		errorHandler(res, req, http.StatusNotFound)
		return
	}
	var nextExercise database.Exercise
	if isOk {
		nextExercise = database.GetExerciseByID(currentLogIn.Progression + 1)
	}
	pageData := homePageVars{
		Title:        "Questline | Home",
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
func errorHandler(res http.ResponseWriter, req *http.Request, status int) {
	if status == http.StatusNotFound {
		tmpl := template.Must(template.ParseFiles("./CLIENT/static/404.gohtml"))
		_ = tmpl.Execute(res, req)
		return
	}
}

func NotLogged(res http.ResponseWriter, req *http.Request) {
	tmpl := template.Must(template.ParseFiles("./CLIENT/static/401.gohtml"))
	_ = tmpl.Execute(res, req)
	return
}

func Favicon(w http.ResponseWriter, r *http.Request) { //Get the favicon route
	http.ServeFile(w, r, "./CLIENT/static/SRC/favicon.png")
}

func ExercicePage(res http.ResponseWriter, req *http.Request) {
	var IsNotHome bool
	var canDo bool = true
	title := strings.Split(req.URL.Path, "/exercice/")[1]
	currentLogIn, isOk := auth.ExtractClaims(res, req)
	if !isOk {
		NotLogged(res, req)
		return
	}
	var currentExersise database.Exercise
	if title != "" {
		currentExersise = database.GetExerciseByName(title)
		IsNotHome = true
	} else {
		currentExersise = database.Exercise{IdExercise: 0, Title: "Welcome to QuestLine", Prompt: "Select an exercise to start learning Go !", Difficulty: 0}
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
	if currentExersise.IdExercise > currentLogIn.Progression+1 {
		canDo = false
	}
	//END PREP
	pageData := exercicePageVars{
		Title:             "QuestLine | " + currentExersise.Title,
		ExerciceTitle:     currentExersise.Title,
		ExercicePrompt:    currentExersise.Prompt,
		ExerciceOutput:    "Click Run Code to test you code !",
		ExercicesList:     exercicesList,
		CanDo:             canDo,
		User:              currentLogIn.Name,
		IsNotHome:         IsNotHome,
		FunctionStructure: currentExersise.FunctionStructure,
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
	Code     string
	Exercice string
}

type Result struct {
	Result string
}

func SendCode(res http.ResponseWriter, req *http.Request) {
	tStart := time.Now()
	var userCode Code
	b, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(b, &userCode)
	if err != nil {
		log.Fatalln(err)
	}
	callingCode := compare.GetTest(userCode.Exercice)
	code := callingCode + "\n" + userCode.Code
	path := container.CreateCodeFile("user", code)
	result := container.TestCode(path)
	if result == nil {
		log.Println("Error compiling")
		_, err := res.Write([]byte("Error"))
		if err != nil {
			return
		}
		return
	}
	err = os.RemoveAll(path)
	if err != nil {
		log.Fatalln(err)
	}
	jsonResult := Result{
		Result: string(result),
	}
	check := compare.Compar(userCode.Exercice, string(result), res, req)
	solution := compare.GetSolution(userCode.Exercice)
	if check {
		jsonResult.Result += "<br><br>&#9989 Well done!<br><br>Code tested in " + strconv.FormatFloat(float64(time.Since(tStart))/1e9, 'f', 2, 64) + " sec"
	} else {
		jsonResult.Result += "<br><br>&#10060 Try again !" + "<br>Expected output : " + solution + "<br><br>Code tested in " + strconv.FormatFloat(float64(time.Since(tStart))/1e9, 'f', 2, 64) + " sec"
	}
	jsonData, err := json.Marshal(jsonResult)
	if err != nil {
		log.Println(err)
	}
	_, err = res.Write(jsonData)
	if err != nil {
		return
	}
}
