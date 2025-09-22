package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderpb "github.com/ojihalaw/shopping-cart-go-grpc/shared/pb/order"
	productpb "github.com/ojihalaw/shopping-cart-go-grpc/shared/pb/product"
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

	rootMux := http.NewServeMux()

	// --- SSE handler ---
	rootMux.HandleFunc("/api/orders/stream", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		ctx := r.Context()
		orderID := r.URL.Query().Get("order_id")

		// Connect ke gRPC OrderService
		conn, err := grpc.NewClient("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			http.Error(w, "failed to connect order service", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		client := orderpb.NewOrderServiceClient(conn)
		stream, err := client.StreamOrderStatus(ctx, &orderpb.OrderStatusRequest{OrderId: orderID})
		if err != nil {
			http.Error(w, "failed to start stream", http.StatusInternalServerError)
			return
		}

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming unsupported", http.StatusInternalServerError)
			return
		}

		for {
			select {
			case <-ctx.Done():
				return
			default:
				res, err := stream.Recv()
				if err != nil {
					if err == io.EOF {
						// ✅ Kirim event completed
						fmt.Fprintf(w, "event: completed\ndata: stream closed\n\n")
						flusher.Flush()
						return
					}
					fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
					flusher.Flush()
					return
				}

				data, _ := json.Marshal(res)
				fmt.Fprintf(w, "data: %s\n\n", data)
				flusher.Flush()
			}
		}
	})

	// --- gRPC-Gateway handler ---
	rootMux.Handle("/", mux)

	// 🔹 Tambahkan CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "https://yourdomain.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(rootMux)

	log.Println("🌐 API Gateway running at :8080")
	if err := http.ListenAndServe(":8080", corsHandler); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
