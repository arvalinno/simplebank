package main

import (
	"context"
	"log"
	"net"

	"github.com/arvalinno/simplebank/api"
	db "github.com/arvalinno/simplebank/db/sqlc"
	"github.com/arvalinno/simplebank/gapi"
	"github.com/arvalinno/simplebank/pb"
	"github.com/arvalinno/simplebank/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot log config:", err)
	}

	connConfig, err := pgx.ParseConfig(config.DBSource)
	if err != nil {
		log.Fatal("cannot parse db config:", err)
	}

	conn, err := pgxpool.New(context.Background(), connConfig.ConnString())
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)

	//if want to run gRPRC Service
	runGrpcServer(config, store)

	// runGinServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("cannot create listener")
	}

	grpcServer := grpc.NewServer()

	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server")
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server")
	}
}
