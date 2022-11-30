package main

import (
	"database/sql"
	"log"

	"github.com/brkss/simplebank/api"
	db "github.com/brkss/simplebank/db/sqlc"
	"github.com/brkss/simplebank/utils"
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/lib/pq"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Invalid Config :", err)
	}
	con, err := sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatal("cannot connect to database : ", err)
	}
	store := db.NewStore(con)
	server := api.NewServer(store)
	err = server.Start(config.ServerAdress)

	if err != nil {
		log.Fatal("cannot connect to server : ", err)
	}
}
