package main

import (
	"context"
	"flag"
	"github.com/mnadev/limestone/internal/infrastructure/server"
	"log"

	"github.com/lpernett/godotenv"
	"github.com/mnadev/limestone/internal/infrastructure/database"
	"github.com/mnadev/limestone/internal/infrastructure/grpc/auth"
)

var (
	grpcEndpoint = flag.String("grpc_endpoint", ":8081", "gRPC server endpoint")
	httpEndpoint = flag.String("http_endpoint", ":8080", "HTTP server endpoint")
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	flag.Parse()
	loadEnv()
	db := database.SetupDatabase()

	// Start gRPC
	grpcServer, grpcListener := server.SetupGRPCServer(db, *grpcEndpoint)
	server.StartGRPCServer(grpcServer, grpcListener)

	// Start REST Gateway
	ctx := context.Background()
	restMux := server.SetupRESTGateway(ctx, db, *grpcEndpoint)
	restHandler := auth.VerifyJWTInterceptorRest(restMux)
	server.StartRESTGateway(ctx, restHandler, *httpEndpoint)
}
