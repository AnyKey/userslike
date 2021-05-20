package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"users_like/grpcsrv/server"
	"users_like/repository"
)

func main() {
	db := initPostgres("dbname=musicdb user=postgres password=123 port=5432 sslmode=disable")
	conn := repository.Repository{Conn: db}
	server.Run(conn)
}

func initPostgres(database string) *sql.DB {
	db, err := sql.Open("postgres", database)
	if err != nil {
		log.Fatalln(err)
	}
	if db.Ping() != nil {
		log.Fatalln(err)
	}
	return db
}
