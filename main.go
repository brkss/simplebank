package main

import (
	"database/sql"
	"log"

	"github.com/brkss/simplebank/api"
	db "github.com/brkss/simplebank/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	DBDRIVER = "postgres"
	DBSOURCE = "postgres://root:root@localhost:5432/simplebank?sslmode=disable"
	ADDRESS  = "0.0.0.0:8080"
)

func main() {

	con, err := sql.Open(DBDRIVER, DBSOURCE)
	if err != nil {
		log.Fatal("cannot connect to database : ", err)
	}
	store := db.NewStore(con)
	server := api.NewServer(store)
	err = server.Start(ADDRESS)

	if err != nil {
		log.Fatal("cannot connect to server : ", err)
	}
}
