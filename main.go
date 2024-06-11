package main

import (
	"context"
	"log"

	"github.com/arvalinno/simplebank/api"
	db "github.com/arvalinno/simplebank/db/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	connConfig, err := pgx.ParseConfig(dbSource)
	if err != nil {
		log.Fatal("cannot parse db config:", err)
	}

	conn, err := pgxpool.New(context.Background(), connConfig.ConnString())
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server")
	}
}
