package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

const (
	DBDRIVER = "postgres"
	DBSOURCE = "postgres://root:root@localhost:5432/simplebank?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(DBDRIVER, DBSOURCE)
	if err != nil {
		log.Fatal("cannot connect to database !", err)
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}
