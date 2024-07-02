package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"restaurant-micro/config"
	"restaurant-micro/jwt"
	restaurantpb "restaurant-micro/proto/restaurant"

	"github.com/gorilla/handlers"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

const (
	StatusBadRequest          = 400
	StatusConflict            = 409
	StatusInternalServerError = 500
	StatusOK                  = 200
	StatusCreated             = 201
	StatusNotFound            = 404
	StatusUnauthorized        = 401
	StatusForbidden           = 403
)

type RestaurantService struct {
	restaurantpb.UnimplementedRestaurantServiceServer
}

var restaurantDBConnector *gorm.DB
var restaurantItemDBConnector *gorm.DB
var restaurantAddressDBConnector *gorm.DB
var logger *zap.Logger

func init() {
	var err error
	logger, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
}

func startServer() {
	if err := godotenv.Load(".env"); err != nil {
		logger.Warn("No .env file found")
	}
	logger.Info("Starting restaurant-microservice server...")
	// Connect to the database
	restaurantDB, restaurantItemDB, restaurantAddress, err := config.ConnectDB()
	restaurantDBConnector = restaurantDB
	restaurantItemDBConnector = restaurantItemDB
	restaurantAddressDBConnector = restaurantAddress

	if err != nil {
		logger.Fatal("Could not connect to the database", zap.Error(err))
	}
	// Start the gRPC server
	listner, err := net.Listen("tcp", "localhost:50052")
	// Check if there is an error while starting the server
	if err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
	// Create a new gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(jwt.UnaryInterceptor),
	)

	// Register the service with the server
	restaurantpb.RegisterRestaurantServiceServer(grpcServer, &RestaurantService{})

	// Start the server in a new goroutine (concurrency) (Serve).
	go func() {
		if err := grpcServer.Serve(listner); err != nil {
			logger.Fatal("Failed to serve", zap.Error(err))
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
		logger.Fatal("Failed to dial server", zap.Error(err))
	}
	// Create a new gRPC-Gateway mux (gateway).
	gwmux := runtime.NewServeMux()

	// Register the service with the server (gateway).
	err = restaurantpb.RegisterRestaurantServiceHandler(context.Background(), gwmux, connection)
	if err != nil {
		logger.Fatal("Failed to register gateway", zap.Error(err))
	}
	// Enable CORS
	corsOrigins := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	corsMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	corsHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	corsHandler := handlers.CORS(corsOrigins, corsMethods, corsHeaders)
	wrappedGwmux := corsHandler(gwmux)
	// Create a new HTTP server (gateway). (Serve). (ListenAndServe)
	gwServer := &http.Server{
		Addr:    ":8091",
		Handler: wrappedGwmux,
	}
	logger.Info("Serving gRPC-Gateway on http://0.0.0.0:8091")
	log.Fatalln(gwServer.ListenAndServe())
}

func main() {
	startServer()
}
