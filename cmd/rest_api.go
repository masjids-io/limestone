package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/mnadev/limestone/gen/go"
	"github.com/mnadev/limestone/internal/application/services"
	"github.com/mnadev/limestone/internal/infrastructure/storage"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func setupRESTGateway(ctx context.Context, db *gorm.DB) *runtime.ServeMux {
	mux := runtime.NewServeMux()
	conn, err := grpc.DialContext(ctx, *grpcEndpoint, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial gRPC server: %v", err)
	}
	defer conn.Close()

	smgr := &storage.StorageManager{DB: db}
	//err = pb.RegisterAdhanServiceHandlerServer(ctx, mux, &adhan_service.AdhanServiceServer{Smgr: smgr})
	//if err != nil {
	//	log.Fatalf("failed to register AdhanService handler: %s", err)
	//}
	//err = pb.RegisterEventServiceHandlerServer(ctx, mux, &event_service.EventServiceServer{Smgr: smgr})
	//if err != nil {
	//	log.Fatalf("failed to register EventService handler: %s", err)
	//}
	//err = pb.RegisterMasjidServiceHandlerServer(ctx, mux, &masjid_service.MasjidServiceServer{Smgr: smgr})
	//if err != nil {
	//	log.Fatalf("failed to register MasjidService handler: %s", err)
	//}
	err = pb.RegisterUserServiceHandlerServer(ctx, mux, &user_service.UserServiceServer{Smgr: smgr})
	if err != nil {
		log.Fatalf("failed to register UserService handler: %s", err)
	}
	return mux
}

func startRESTGateway(ctx context.Context, handler http.Handler) {
	log.Printf("HTTP server listening on %s", *httpEndpoint)
	if err := http.ListenAndServe(*httpEndpoint, handler); err != nil {
		log.Fatalf("failed to serve HTTP traffic: %s", err)
	}
}
