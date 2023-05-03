package database

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	IdUser      int
	Progression int
	Name        string
	Email       string
	Pwd         string
}

type Exercise struct {
	IdExercise int
	Title      string
	Prompt     string
	Difficulty int
}

// ---INIT---

// Singleton for db instance
var dbInstance *sql.DB

func OpenDb() (*sql.DB, error) {
	file := "./DB/database.db"
	sqlDB, connectionError := sql.Open("sqlite3", file)
	return sqlDB, connectionError
}

func GetDbInstance() *sql.DB {
	if dbInstance == nil {
		var ConnectionErr error
		dbInstance, ConnectionErr = OpenDb()
		if ConnectionErr != nil {
			log.Println("[DATABASE] Failed to connect database [ERR]: ", ConnectionErr)
		} else {
			log.Println("[DATABASE] Connecting database...")
		}
	}
	return dbInstance
}

// ---AUTH---

func CheckUserExist(email string) bool {
	var count int
	err := GetDbInstance().QueryRow("SELECT count(*) FROM User WHERE Email = ?", email).Scan(&count)
	if err != nil {
		log.Println("[DATABASE] Failed to check user exist [ERR]: ", err)
		return true
	}
	if count == 0 {
		return false
	}
	return true
}

func InsertUser(name string, email string, password string) {
	password = HashPassword(password)
	_, err := GetDbInstance().Exec("INSERT INTO User (Progression, Username, Email, Password, Exp) VALUES (?, ?, ?, ?, ?)", 0, name, email, password, 0)
	if err != nil {
		log.Println("[DATABASE] Failed to insert user [ERR]: ", err)
	}
}

func GetUserPassword(email string) string {
	var password string
	err := GetDbInstance().QueryRow("SELECT Password FROM User WHERE Email = ?", email).Scan(&password)
	if err != nil {
		log.Println("[DATABASE] Failed to get user password [ERR]: ", err)
		return ""
	}
	return password
}

// ---UTILS---

func HashPassword(password string) string {
	h := sha1.New()
	h.Write([]byte(password))
	sha1Hash := hex.EncodeToString(h.Sum(nil))
	return sha1Hash
}

func GetUserByMail(mail string) User {
	var user User
	err := GetDbInstance().QueryRow("SELECT * FROM User WHERE Email = ?", mail).Scan(&user.IdUser, &user.Progression, &user.Name, &user.Email, &user.Pwd)
	if err != nil {
		log.Println("[DATABASE] Failed to get user by mail [ERR]: ", err)

	}
	return user
}

func GetExerciseNameList() []string {
	var exerciceNameList []string
	rows, err := GetDbInstance().Query("SELECT Title FROM Exercise")
	if err != nil {
		log.Println("[DATABASE] Failed to get exercice name list [ERR]: ", err)
	}
	defer rows.Close()
	for rows.Next() {
		var exerciceName string
		err := rows.Scan(&exerciceName)
		if err != nil {
			log.Println("[DATABASE] Failed to get exercice name list [ERR]: ", err)
		}
		exerciceNameList = append(exerciceNameList, exerciceName)
	}
	return exerciceNameList
}

func GetExerciseByName(name string) Exercise {
	var exercise Exercise
	err := GetDbInstance().QueryRow("SELECT * FROM Exercise WHERE Title = ?", name).Scan(&exercise.IdExercise, &exercise.Title, &exercise.Prompt, &exercise.Difficulty)
	if err != nil {
		log.Println("[DATABASE] Failed to get exercise by name [ERR]: ", err)
		return Exercise{}
	}
	return exercise
}

func GetExerciseByID(ID int) Exercise {
	var exercise Exercise
	err := GetDbInstance().QueryRow("SELECT * FROM Exercise WHERE ID_Exercise = ?", ID).Scan(&exercise.IdExercise, &exercise.Title, &exercise.Prompt, &exercise.Difficulty)
	if err != nil {
		log.Println("[DATABASE] Failed to get exercise by ID [ERR]: ", err)
		return Exercise{}
	}
	return exercise
}
