package database

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type User struct {
	IdUser      int
	Progression int
	Name        string
	Email       string
	Pwd         string
	Score       int
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
	err := GetDbInstance().QueryRow("SELECT count(*) FROM Users WHERE Email = ?", email).Scan(&count)
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
	_, err := GetDbInstance().Exec("INSERT INTO Users (Progression, Name, Email, Pwd, Score) VALUES (?, ?, ?, ?, ?)", 0, name, email, password, 0)
	if err != nil {
		log.Println("[DATABASE] Failed to insert user [ERR]: ", err)
	}
}

func GetUserPassword(email string) string {
	var password string
	err := GetDbInstance().QueryRow("SELECT Pwd FROM Users WHERE Email = ?", email).Scan(&password)
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
	err := GetDbInstance().QueryRow("SELECT * FROM Users WHERE Email = ?", mail).Scan(&user.IdUser, &user.Progression, &user.Name, &user.Email, &user.Pwd, &user.Score)
	if err != nil {
		log.Println("[DATABASE] Failed to get user by mail [ERR]: ", err)

	}
	return user
}
