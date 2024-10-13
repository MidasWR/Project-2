package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

const (
	host     = "localhost"
	user     = "midas"
	password = "123qwer321QWER"
	dbname   = "users"
	port     = "5432"
)

func NewDB() *sql.DB {
	con := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", con)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		username varchar(20),
		login varchar(30),
		password varchar(100)) ;
	`)
	if err != nil {
		panic(err)
	}
	log.Println("DB creating access")
	return db
}
