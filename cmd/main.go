package main

import (
	"context"
	"flag"
	"github.com/mnadev/limestone/internal/infrastructure/auth"
	"github.com/mnadev/limestone/internal/infrastructure/server"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/lpernett/godotenv"
	"github.com/mnadev/limestone/internal/infrastructure/database"
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
	//ctx := context.Background()
	//restMux := server.SetupRESTGateway(ctx, db)
	//restHandler := auth.VerifyJWTInterceptorRest(restMux)
	//server.StartRESTGateway(restHandler, *httpEndpoint)

	mainMux := http.NewServeMux()

	ctx := context.Background()
	grpcGatewayMux := server.SetupRESTGateway(ctx, db)

	restHandlerWithAuth := auth.VerifyJWTInterceptorRest(grpcGatewayMux)

	mainMux.Handle("/", restHandlerWithAuth)

	mainMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	staticFilesPath := filepath.Join(currentDir, "docs")

	if _, err := os.Stat(staticFilesPath); os.IsNotExist(err) {
		log.Printf("Warning: Static files directory '%s' does not exist. Ensure 'static' folder is copied to /app/static in Dockerfile.", staticFilesPath)
	} else {
		log.Printf("Serving static files from: %s", staticFilesPath)
	}

	fs := http.FileServer(http.Dir(staticFilesPath))
	mainMux.Handle("/docs/", http.StripPrefix("/docs/", fs))

	mainMux.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
	})

	log.Printf("HTTP Server listening on %s", *httpEndpoint)
	log.Fatal(http.ListenAndServe(*httpEndpoint, mainMux))
}
