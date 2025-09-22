package main

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderpb "github.com/ojihalaw/sample-grpc/shared/pb/order"
	productpb "github.com/ojihalaw/sample-grpc/shared/pb/product"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// 🔹 Register ProductService
	if err := productpb.RegisterProductServiceHandlerFromEndpoint(
		ctx, mux, "localhost:50052", opts,
	); err != nil {
		log.Fatalf("failed to register ProductService: %v", err)
	}

	// 🔹 Register OrderService
	if err := orderpb.RegisterOrderServiceHandlerFromEndpoint(
		ctx, mux, "localhost:50053", opts,
	); err != nil {
		log.Fatalf("failed to register OrderService: %v", err)
	}

	log.Println("🌐 API Gateway running at :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
