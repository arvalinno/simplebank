package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/arvalinno/simplebank/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	connConfig, err := pgx.ParseConfig(config.DBSource)
	if err != nil {
		log.Fatal("cannot parse db config:", err)
	}

	testDB, err = pgxpool.New(context.Background(), connConfig.ConnString())
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
