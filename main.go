package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/giadat1599/small_bank/api"
	db "github.com/giadat1599/small_bank/db/sqlc"
	"github.com/giadat1599/small_bank/gapi"
	"github.com/giadat1599/small_bank/pb"
	"github.com/giadat1599/small_bank/utils"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	runGRPCServer(config, store)

}

func runGRPCServer(config utils.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot create the server: ", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSmallBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddr)

	if err != nil {
		log.Fatal("cannot create listener: ", err)
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start the grpc server: ", err)
	}
}

func runGinServer(config utils.Config, store db.Store) {

	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal("Cannot create the server: ", err)
	}

	err = server.StartServer(config.HTTPServerAddr)

	if err != nil {
		log.Fatal("Cannot start the server: ", err)
	}
}