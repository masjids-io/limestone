package main

import (
	"context"
	"flag"
	"github.com/lpernett/godotenv"
	"github.com/mnadev/limestone/internal/infrastructure/database"
	"github.com/mnadev/limestone/internal/infrastructure/grpc/auth"
	"log"
)

var (
	grpcEndpoint = flag.String("grpc_endpoint", "localhost:8081", "gRPC server endpoint")
	httpEndpoint = flag.String("http_endpoint", "localhost:8080", "HTTP server endpoint")
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	loadEnv()
	db := database.SetupDatabase()
	grpcServer, grpcListener := setupGRPCServer(db)
	startGRPCServer(grpcServer, grpcListener)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	restMux := setupRESTGateway(ctx, db)
	restHandler := auth.VerifyJWTInterceptorRest(restMux)
	startRESTGateway(ctx, restHandler)
}
