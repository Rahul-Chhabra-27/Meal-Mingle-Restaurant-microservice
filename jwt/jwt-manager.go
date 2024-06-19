package jwt

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}
type UserClaims struct {
	jwt.StandardClaims
	UserEmail string 
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) (*JWTManager, error) {
	return &JWTManager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}, nil
}

func VerifyToken(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}
			return []byte(os.Getenv("SECRET_KEY")), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
// Middleware to check if the user is authenticated or not by checking the JWT token provided in the request
func UnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	// skip the authentication for the health check endpoint for GetAllRestaurants
	if info.FullMethod == "/proto.RestaurantService/GetRestaurantsByCity" {
		return handler(ctx, req)
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}
	tokenString := md.Get("authorization")
	if len(tokenString) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}
	token := strings.Split(tokenString[0], " ")
	// Parse JWT token
	claims, err := VerifyToken(token[1])
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "token is invalid: %v", err)
	}
	// Pass useremail to context for further use
	ctx = context.WithValue(ctx, "userEmail", claims.UserEmail)
	// Proceed with the request
	return handler(ctx, req)
}
