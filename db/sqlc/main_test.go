package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/brkss/simplebank/utils"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config : ", err)
	}
	testDB, err = sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatal("cannot connect to database !", err)
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}
