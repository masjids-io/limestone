package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/go-redis/redis"
	"github.com/gorilla/sessions"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"

	"github.com/mnadev/limestone/adhan_service"
	"github.com/mnadev/limestone/auth"
	"github.com/mnadev/limestone/event_service"
	"github.com/mnadev/limestone/masjid_service"
	pb "github.com/mnadev/limestone/proto"
	"github.com/mnadev/limestone/storage"
	"github.com/mnadev/limestone/user_service"
)

var (
	grpcEndpoint = flag.String("grpc_endpoint", ":8081", "gRPC server endpoint")
	httpEndpoint = flag.String("http_endpoint", ":8080", "HTTP server endpoint")
)

func main() {
	lis, err := net.Listen("tcp", *grpcEndpoint)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Include AuthMiddleware as a UnaryInterceptor
	server := grpc.NewServer(
		grpc.UnaryInterceptor(auth.AuthMiddleware),
	)

	host := os.Getenv("db-host")
	port := os.Getenv("db-port")
	dbName := os.Getenv("db-name")
	dbUser := os.Getenv("db-user")
	password := os.Getenv("db-password")
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host,
		port,
		dbUser,
		dbName,
		password,
	)
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	DB.AutoMigrate(storage.AdhanFile{})
	DB.AutoMigrate(storage.Event{})
	DB.AutoMigrate(storage.Masjid{})
	DB.AutoMigrate(storage.User{})

	// redis connection has been created
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redis-url"),
		Password: os.Getenv("redis-pass"), // password set
		DB:       0,                       // use default DB
	})
	res := rdb.Ping()
	fmt.Println("Redis response: " + res.String()) // response being pinged
	adhan_service_server := adhan_service.AdhanServiceServer{
		SM: &storage.StorageManager{
			DB:    DB,
			Cache: rdb,
		},
	}
	gothic.Store = sessions.NewCookieStore([]byte(os.Getenv("JWT_KEY")))
	goth.UseProviders(
		google.New(os.Getenv("google-client"), os.Getenv("google-secret"), os.Getenv("google-scope")),
		//TODO: ADD MICROSOFT SUPPORT
		//microsoftonline.New(os.Getenv("microsoft-client"), os.Getenv("microsoft-secret"), os.Getenv("microsoft-scope")),
	)
	pb.RegisterAdhanServiceServer(server, &adhan_service_server)
	event_server := event_service.EventServiceServer{
		SM: &storage.StorageManager{
			DB:    DB,
			Cache: rdb,
		},
	}
	pb.RegisterEventServiceServer(server, &event_server)
	masjid_server := masjid_service.MasjidServiceServer{
		SM: &storage.StorageManager{
			DB:    DB,
			Cache: rdb,
		},
	}
	pb.RegisterMasjidServiceServer(server, &masjid_server)
	user_server := user_service.UserServiceServer{
		SM: &storage.StorageManager{
			DB:    DB,
			Cache: rdb,
		},
	}
	pb.RegisterUserServiceServer(server, &user_server)

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC traffic: %s", err)
		}
	}()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	// Create a new http.ServeMux for your custom SSO routes
	customMux := http.NewServeMux()
	httpMux := http.NewServeMux()
	err = pb.RegisterAdhanServiceHandlerServer(ctx, mux, &adhan_service_server)
	if err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
	err = pb.RegisterEventServiceHandlerServer(ctx, mux, &event_server)
	if err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
	err = pb.RegisterMasjidServiceHandlerServer(ctx, mux, &masjid_server)
	if err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
	err = pb.RegisterUserServiceHandlerServer(ctx, mux, &user_server)
	if err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
	app := App{
		DB: DB,
	}
	// Register your OAuth routes
	customMux.HandleFunc("/auth/google", app.HandleGoogleOauthRoute)
	customMux.HandleFunc("/auth/google/callback", app.HandleGoogleOauthCallbackRoute)
	// Use the main mux to delegate requests
	httpMux.Handle("/", auth.HttpAuthMiddleware(mux))  // Handle gRPC-Gateway routes
	httpMux.Handle("/auth/google", customMux)          // Handle custom http routes for google
	httpMux.Handle("/auth/google/callback", customMux) //Handle custom http google callback route
	if err = http.ListenAndServe(*httpEndpoint, httpMux); err != nil {
		log.Fatalf("failed to serve HTTP traffic: %s", err)
	}
}
