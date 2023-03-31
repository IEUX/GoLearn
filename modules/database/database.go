package database

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"log"
)

var dbInstance *sql.DB

//INIT

func OpenDb() (*sql.DB, error) {
	file := "database.db"
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

func HashPassword(password string) string {
	h := sha1.New()
	h.Write([]byte(password))
	sha1Hash := hex.EncodeToString(h.Sum(nil))
	return sha1Hash
}
