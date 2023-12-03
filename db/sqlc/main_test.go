package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Use local postgres to run tests
const (
	dbSource = "postgres://root:secret@localhost:5432/small_bank?sslmode=disable"
)

var testStore Store

func TestMain(m *testing.M) {
	var err error
	connPool, err := pgxpool.New(context.Background(), dbSource)

	if err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}

	testStore = NewStore(connPool)

	os.Exit(m.Run())
}
