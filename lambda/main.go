package main

import (
	"context"
	"fmt"
	"zombie-boss-api/pkg/database"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.Path {
	case "/zombie-boss/count/set/{count}":
		return handleSetCount(ctx, request)
	case "/zombie-boss/count/add":
		return handleAddCount(ctx, request)
	case "/zombie-boss/count":
		return handleGetCount(ctx, request)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Not Found: " + request.Path,
		}, nil
	}
}

func handleSetCount(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	count := request.PathParameters["count"]

	err := database.SetCount(count)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error setting count: %v", err),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Count set successfully",
	}, nil
}

func handleAddCount(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	err := database.AddCount()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error adding to count: %v", err),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Count added successfully",
	}, nil
}

func handleGetCount(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	count, err := database.GetCount()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error getting count: %v", err),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("Current count is: %v", count),
	}, nil
}

func main() {
	lambda.Start(handler)
}
