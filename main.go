package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	db "github.com/giadat1599/small_bank/db/sqlc"
	_ "github.com/giadat1599/small_bank/doc/statik"
	"github.com/giadat1599/small_bank/gapi"
	"github.com/giadat1599/small_bank/pb"
	"github.com/giadat1599/small_bank/utils"
	"github.com/giadat1599/small_bank/worker"
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
		log.Fatal().Msg("Cannot load configuration")
	}
	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	connection, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal().Msg("Cannot connect to database")
	}

	runDBMirgation(config.MigrationURL, config.DBSource)

	store := db.NewStore(connection)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddr,
	}

	taskDistributor := worker.NewRedisDistributor(redisOpt)

	go runTaskProcessor(redisOpt, store)
	go runGatewayServer(config, store, taskDistributor)
	runGRPCServer(config, store, taskDistributor)

}

func runDBMirgation(migrationURL string, dbSource string) {
	m, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Msg("cannot create new migrate instance")
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Msg("failed to run migrate up")
	}

	log.Info().Msg("db migrated succesfully")
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)
	log.Info().Msg("start task processor")
	err := taskProcessor.Start()

	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}

func runGatewayServer(config utils.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msg("Cannot create the server")
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
		log.Fatal().Msg("cannot register handler server")
	}

	// Receive http requests from client
	mux := http.NewServeMux()
	// Convert to grpc format
	mux.Handle("/", grpcMux)
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Msg("cannot create statik fs")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddr)

	if err != nil {
		log.Fatal().Msg("cannot create listener")
	}

	log.Info().Msgf("start HTTP server at %s", listener.Addr().String())
	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Msg("cannot start the grpc server")
	}
}

/* Serving gRPC */
func runGRPCServer(config utils.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msg("Cannot create the server")
	}
	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSmallBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddr)

	if err != nil {
		log.Fatal().Msg("cannot create listener")
	}

	log.Info().Msgf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Msg("cannot start the grpc server")
	}
}

/* Serving HTTP requests using Gin framework */
// func runGinServer(config utils.Config, store db.Store) {

// 	server, err := api.NewServer(config, store)

// 	if err != nil {
// 		log.Fatal().Msg("Cannot create the server")
// 	}

// 	err = server.StartServer(config.HTTPServerAddr)

// 	if err != nil {
// 		log.Fatal().Msg("Cannot start the server")
// 	}
// }
