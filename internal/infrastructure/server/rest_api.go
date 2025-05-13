package server

import (
	"context"
	pb "github.com/mnadev/limestone/gen/go"
	services "github.com/mnadev/limestone/internal/application/services"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/mnadev/limestone/internal/application/handler"
	"github.com/mnadev/limestone/internal/infrastructure/storage"
	"gorm.io/gorm"
)

func SetupRESTGateway(ctx context.Context, db *gorm.DB, grpcEndpoint string) *runtime.ServeMux {
	mux := runtime.NewServeMux()

	repo := storage.NewGormUserRepository(db)
	svc := services.NewUserService(repo)
	handler := handler.NewUserGrpcHandler(svc)

	if err := pb.RegisterUserServiceHandlerServer(ctx, mux, handler); err != nil {
		log.Fatalf("failed to register UserService handler: %s", err)
	}
	return mux
}

func StartRESTGateway(ctx context.Context, handler http.Handler, httpEndpoint string) {
	log.Printf("HTTP server listening on %s", httpEndpoint)
	if err := http.ListenAndServe(httpEndpoint, handler); err != nil {
		log.Fatalf("failed to serve HTTP traffic: %s", err)
	}
}
