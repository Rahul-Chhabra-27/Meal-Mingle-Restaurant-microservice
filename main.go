package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	restaurantpb "restaurant-micro/proto/restaurant"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RestaurantService struct {
	restaurantpb.UnimplementedRestaurantServiceServer
}

func startServer() {
	fmt.Println("Starting restaurant-microservice server...")

	// Start the gRPC server
	listner, err := net.Listen("tcp", "localhost:50052")
	// Check if there is an error while starting the server
	if err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the service with the server
	restaurantpb.RegisterRestaurantServiceServer(grpcServer, &RestaurantService{})
	
	// Start the server in a new goroutine (concurrency) (Serve).
	go func() {
		if err := grpcServer.Serve(listner); err != nil {
			log.Fatalf("Failed to serve: %s", err)
		}
	}()
	// Create a new gRPC-Gateway server (gateway).
	connection, err := grpc.DialContext(
		context.Background(),
		"localhost:50052",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	// Create a new gRPC-Gateway mux (gateway).
	gwmux := runtime.NewServeMux()

	// Register the service with the server (gateway).
	err = restaurantpb.RegisterRestaurantServiceHandler(context.Background(), gwmux, connection)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	// Create a new HTTP server (gateway). (Serve). (ListenAndServe)
	gwServer := &http.Server{
		Addr:    ":8091",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8091")
	log.Fatalln(gwServer.ListenAndServe())
}

func main() {
	startServer()
}
