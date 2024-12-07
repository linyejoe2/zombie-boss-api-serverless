package main

import (
	"context"
	"fmt"
	"strings"
	"zombie-boss-api/pkg/database"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var headers = map[string]string{
	"Access-Control-Allow-Origin":  "*",
	"Access-Control-Allow-Methods": "GET, POST, OPTIONS",
	"Access-Control-Allow-Headers": "Content-Type",
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if strings.Contains(request.Path, "/count/set") {
		return handleSetCount(ctx, request)
	} else if strings.Contains(request.Path, "/count/add") {
		return handleAddCount(ctx, request)
	} else if strings.Contains(request.Path, "/count") {
		return handleGetCount(ctx, request)
	} else if strings.Contains(request.Path, "/pre-register") {
		return handlePreRegister(ctx, request)
	} else {
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			StatusCode: 404,
			Body:       "Not Found: " + request.Path,
		}, nil
	}
}

func handlePreRegister(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	email := request.PathParameters["email"]

	err := database.PreRegister(email)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			StatusCode: 500,
			Body:       fmt.Sprintf(`{"error": true, "message": "Error registering: %v"}`, err),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Headers:    headers,
		StatusCode: 200,
		Body:       `{"error": false, "message": "Successfully registered"}`,
	}, nil
}

func handleSetCount(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	count := request.PathParameters["count"]

	err := database.SetCount(count)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			StatusCode: 500,
			Body:       fmt.Sprintf("Error setting count: %v", err),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Headers:    headers,
		StatusCode: 200,
		Body:       "Count set successfully",
	}, nil
}

func handleAddCount(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	err := database.AddCount()
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			StatusCode: 500,
			Body:       fmt.Sprintf("Error adding to count: %v", err),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Headers:    headers,
		StatusCode: 200,
		Body:       "Count added successfully",
	}, nil
}

func handleGetCount(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	count, err := database.GetCount()
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			StatusCode: 500,
			Body:       fmt.Sprintf("Error getting count: %v", err),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Headers:    headers,
		StatusCode: 200,
		Body:       fmt.Sprintf("%v", count),
	}, nil
}

func main() {
	lambda.Start(handler)
}
