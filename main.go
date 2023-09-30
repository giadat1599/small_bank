package main

import (
	"database/sql"
	"log"

	"github.com/giadat1599/small_bank/api"
	db "github.com/giadat1599/small_bank/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver   = "postgres"
	dbSource   = "postgres://root:secret@localhost:5432/small_bank?sslmode=disable"
	serverAddr = "0.0.0.0:8080"
)

func main() {
	connection, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}

	store := db.NewStore(connection)

	server := api.NewServer(store)

	err = server.StartServer(serverAddr)

	if err != nil {
		log.Fatal("Cannot start the server: ", err)
	}

}
