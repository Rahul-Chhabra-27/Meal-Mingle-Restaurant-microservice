package main

import (
	"context"
	"fmt"
	restaurantpb "restaurant-micro/proto/restaurant"

	"google.golang.org/grpc/codes"
)

func (*RestaurantService) DeleteRestaurantItem(ctx context.Context, response *restaurantpb.DeleteRestaurantItemRequest) (*restaurantpb.DeleteRestaurantItemResponse, error) {
	// Get the user email from the context
	userEmail, ok := ctx.Value("userEmail").(string)
	if !ok {
		fmt.Println("Failed to get user email from context")
		return &restaurantpb.DeleteRestaurantItemResponse{
			Message:    "",
			StatusCode: int64(codes.Internal),
			Error:      "Internal Server Error",
		}, nil
	}
	fmt.Println(userEmail);
	return nil,nil;
}
