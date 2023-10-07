package main

import (
	"database/sql"
	"log"

	"github.com/giadat1599/small_bank/api"
	db "github.com/giadat1599/small_bank/db/sqlc"
	"github.com/giadat1599/small_bank/utils"
	_ "github.com/lib/pq"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load configuration")
	}
	connection, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}

	store := db.NewStore(connection)

	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal("Cannot create the server: ", err)
	}

	err = server.StartServer(config.ServerAddr)

	if err != nil {
		log.Fatal("Cannot start the server: ", err)
	}

}
