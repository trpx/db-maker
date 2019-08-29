package utils

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type DB struct {
	User string
	Password string
	Name string
	Host string
	Port string

	backend *sql.DB
}

func (db *DB) Connect() {
	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		db.User,
		db.Password,
		db.Name,
		db.Host,
		db.Port,
	)
	var err error
	db.backend, err = sql.Open("postgres", connStr)
	if err != nil {
		Panicf("couldn't connect to db '%s':\n%v", db.Name, err)
	}
}

// defer this right after .Connect()
func (db *DB) Disconnect() {
	err := db.backend.Close()
	if err != nil {
		Panicf("couldn't close connection to db '%s':\n%v", db.Name, err)
	}
}

// todo: create user, password-related stuff etc

func (db *DB) CreateDBIfNotExists(dbname string) (created bool) {
	exists := db.DbExists(dbname)

	if ! exists {
		// seems like queries with databases like create / drop db etc don't support $1, $2, ... params
		_, err := db.backend.Exec(fmt.Sprintf("CREATE database %s", dbname))
		if err != nil {
			Panicf("couldn't create database '%s'\n%v", dbname, err)
		}
		created = true
	} else {
		// already exists
		created = false
	}
	return created
}

func (db *DB) DbExists (dbname string) (exists bool) {
	var dbCount int
	err := db.backend.QueryRow("SELECT count(*) FROM pg_database WHERE datname = $1", dbname).Scan(&dbCount)
	if err != nil {
		Panicf("couldn't fetch existence of db '%s':\n%v", dbname, err)
	}
	exists = dbCount > 0
	return exists
}
