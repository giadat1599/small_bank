package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	db "github.com/giadat1599/small_bank/db/sqlc"
	_ "github.com/giadat1599/small_bank/doc/statik"
	"github.com/giadat1599/small_bank/gapi"
	"github.com/giadat1599/small_bank/pb"
	"github.com/giadat1599/small_bank/utils"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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

	runDBMirgation(config.MigrationURL, config.DBSource)

	store := db.NewStore(connection)
	// We need to run the gateway or gRPC server in a separate go routine than the main routine to avoid blocking each other
	go runGatewayServer(config, store)
	runGRPCServer(config, store)

}

func runDBMirgation(migrationURL string, dbSource string) {
	m, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create new migrate instance: ", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up: ", err)
	}

	log.Println("db migrated succesfully")
}

func runGatewayServer(config utils.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot create the server: ", err)
	}
	
	// use the exact names which are defined in the proto files for http requests/responses
	jsonOptions := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	
	grpcMux := runtime.NewServeMux(jsonOptions)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterSmallBankHandlerServer(ctx, grpcMux, server)

	if err != nil {
		log.Fatal("cannot register handler server: ", err)
	}

	// Receive http requests from client
	mux := http.NewServeMux()
	// Convert to grpc format
	mux.Handle("/", grpcMux)
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal("cannot create statik fs: ", err)
	}

	swaggerHandler :=  http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddr)

	if err != nil {
		log.Fatal("cannot create listener: ", err)
	}

	log.Printf("start HTTP server at %s", listener.Addr().String())

	err = http.Serve(listener,mux)
	if err != nil {
		log.Fatal("cannot start the grpc server: ", err)
	}
}

/* Serving gRPC */
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

/* Serving HTTP requests using Gin framework */
// func runGinServer(config utils.Config, store db.Store) {

// 	server, err := api.NewServer(config, store)

// 	if err != nil {
// 		log.Fatal("Cannot create the server: ", err)
// 	}

// 	err = server.StartServer(config.HTTPServerAddr)

// 	if err != nil {
// 		log.Fatal("Cannot start the server: ", err)
// 	}
// }