package utils

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	_ "github.com/trpx/pq"
	"io"
	"strings"
)

type DB struct {
	User      string
	Passwords []string
	Name      string
	Host      string
	Port      string

	backend *sql.DB
}

func (db *DB) Connect() {
	var err error
	for _, pw := range db.Passwords {
		connStr := fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
			db.User,
			pw,
			db.Name,
			db.Host,
			db.Port,
		)
		db.backend, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.backend.Ping()
			if err == nil {
				break
			}
		}

	}
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

func (db *DB) CreateDBIfNotExistsWithOwner(dbname string, ownerUsername string) (created bool) {
	exists := db.DBExists(dbname)
	if exists {
		return false
	}

	// seems like queries with databases like create / drop db etc don't support $1, $2, ... params
	_, err := db.backend.Exec(fmt.Sprintf("CREATE DATABASE %s WITH OWNER %s", dbname, ownerUsername))
	if err != nil {
		Panicf("couldn't create database '%s' with owner '%s'\n%v", dbname, ownerUsername, err)
	}

	return true
}

func (db *DB) CreateExtensionIfNotExists(extension string) (created bool) {
	extensionExists := db.ExtensionExists(extension)
	if extensionExists {
		return false
	}
	_, err := db.backend.Exec(fmt.Sprintf("CREATE EXTENSION %s", extension))
	if err != nil {
		Panicf("couldn't create extension '%s'\n%v", extension, err)
	}
	return true
}

func (db *DB) ExtensionExists(extension string) (exists bool) {
	var extCount int
	err := db.backend.QueryRow("SELECT count(*) FROM pg_extension WHERE extname=$1", extension).Scan(&extCount)
	if err != nil {
		Panicf("couldn't fetch existence of extension '%s':\n%v", extension, err)
	}
	exists = extCount > 0
	return exists
}

func (db *DB) DBExists(dbname string) (exists bool) {
	var dbCount int
	err := db.backend.QueryRow("SELECT count(*) FROM pg_database WHERE datname = $1", dbname).Scan(&dbCount)
	if err != nil {
		Panicf("couldn't fetch existence of db '%s':\n%v", dbname, err)
	}
	exists = dbCount > 0
	return exists
}

func (db *DB) UpdateUserPassword(username string, newPassword string) (updated bool) {
	userPasswordHash := db.FetchUserPasswordHash(username)
	newPasswordHash := hashUserPassword(username, newPassword)
	if userPasswordHash == newPasswordHash {
		return false
	}
	_, err := db.backend.Exec(fmt.Sprintf("ALTER ROLE %s WITH ENCRYPTED PASSWORD '%s'", username, newPasswordHash))
	if err != nil {
		Panicf("couldn't change user '%s' password", username)
	}
	return true
}

func (db *DB) FetchUserPasswordHash(username string) (passwordHash string) {
	row := db.backend.QueryRow("SELECT rolpassword FROM pg_authid WHERE rolname=$1", username)
	err := row.Scan(&passwordHash)
	if err != nil {
		print(strings.Contains(err.Error(), "password"))
		Panicf("couldn't fetch password hash of user '%s':\n%v", username, err)
	}
	return passwordHash
}

func (db *DB) CreateUserIfNotExists(username string, password string) (created bool) {
	userExists := db.UserExists(username)
	if userExists {
		return false
	}

	hashedPassword := hashUserPassword(username, password)

	createUserSql := fmt.Sprintf("CREATE ROLE %s WITH ENCRYPTED PASSWORD '%s' LOGIN", username, hashedPassword)
	_, err := db.backend.Exec(createUserSql)
	if err != nil {
		Panicf("couldn't create user '%s':\n%v", username, err)
	}
	return true
}

func hashUserPassword(username string, password string) (hashedPassword string) {
	h := md5.New()
	_, err := io.WriteString(h, password+username)
	if err != nil {
		Panicf("couldn't hash user's '%s' password+username", username)
	}
	hashedPassword = fmt.Sprintf("md5%x", h.Sum(nil))
	return hashedPassword
}

func (db *DB) UserExists(username string) (exists bool) {
	var userCount int
	err := db.backend.QueryRow("SELECT count(*) FROM pg_roles WHERE rolname=$1", username).Scan(&userCount)
	if err != nil {
		Panicf("couldn't fetch existence of user '%s':\n%v", username, err)
	}
	exists = userCount > 0
	return exists
}
