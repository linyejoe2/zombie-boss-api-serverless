package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	utils "zombie-boss-api/pkg"
	"zombie-boss-api/pkg/database"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var headers = map[string]string{
	"Access-Control-Allow-Origin":  "*",
	"Access-Control-Allow-Methods": "GET, POST, OPTIONS",
	"Access-Control-Allow-Headers": "Content-Type",
}

var ginLambda *ginadapter.GinLambdaV2

func init() {
	// 初始化 Gin
	r := gin.Default()

	r.Use(corsMiddleware())

	//! test connection
	// r.Use(func(ctx *gin.Context) {
	// 	path := ctx.FullPath()
	// 	log.Printf(fmt.Sprintf("path: %v", path))
	// 	ctx.JSON(200, path)
	// })

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/count", func(ctx *gin.Context) {
		count, err := database.GetCount()
		if err != nil {
			utils.JSONErrorResponse(ctx, http.StatusInternalServerError, "getting count", err)
			// ctx.JSON(500, gin.H{"message": fmt.Sprintf("Error getting count: %v", err)})
			return
		}

		utils.JSONResponse(ctx, http.StatusOK, false, "Getting count Successfully.", gin.H{"count": count})
		// ctx.JSON(http.StatusOK, gin.H{"body": gin.H{"count": count}})
	})

	r.POST("/count/:count", func(ctx *gin.Context) {
		count := ctx.Param("count")
		err := database.SetCount(count)
		if err != nil {
			utils.JSONErrorResponse(ctx, 500, "setting count", err)
			// ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Error setting count: %v", err)})
		}

		utils.JSONResponse(ctx, 200, false, "Setting count successfully.", nil)
		// ctx.JSON(http.StatusOK, gin.H{"message": "Count set successfully."})
	})

	r.PUT("/count", func(ctx *gin.Context) {
		err := database.AddCount()
		if err != nil {
			utils.JSONErrorResponse(ctx, 500, "increasing count", err)
			// ctx.JSON(http.StatusInternalServerError, utils.JSONResponse(true, fmt.Sprintf("Failed to add count: %v", err), ""))
			return
		}

		utils.JSONSuccessResponse(ctx, "Increasing count", nil)
		// ctx.JSON(200, gin.H{"message": "Count added successfully."})
	})

	r.POST("/pre-register/:email", func(ctx *gin.Context) {
		email := ctx.Param("email")
		emailRegex, emailRegErr := regexp.Compile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
		if emailRegErr != nil {
			utils.JSONErrorResponse(ctx, 500, "inserting pre-register email", emailRegErr)
			return
		}
		if !emailRegex.MatchString(email) {
			utils.JSONErrorResponse(ctx, 400, "inserting pre-register email", errors.New("Unsupported prameters."))
			return
		}

		err := database.PreRegister(email)
		if err != nil {
			utils.JSONErrorResponse(ctx, 500, "inserting pre-register email", err)
			// ctx.JSON(500, gin.H{"message": fmt.Sprintf(`"message": "Error registering: %v"}`, err)})
			return
		}

		utils.JSONSuccessResponse(ctx, "Inserting pre-register email", nil)
		// ctx.JSON(200, gin.H{"message": "Count added successfully."})
	})

	r.GET("/pre-register", func(ctx *gin.Context) {
		emails, err := database.GetPreRegister()
		if err != nil {
			utils.JSONErrorResponse(ctx, 500, "getting pre-register list", err)
			// ctx.JSON(500, gin.H{"message": fmt.Sprintf(`"message": "Error getting register list: %v"}`, err)})
			return
		}

		utils.JSONSuccessResponse(ctx, "Getting pre-register list", emails)
		// ctx.JSON(200, utils.JSONResponse(false, 200, "Pre register list get successfully.", emails))
	})

	r.GET("/register", func(c *gin.Context) {
		registerList, err := database.GetRegister()
		if err != nil {
			utils.JSONErrorResponse(c, 500, "getting register list", err)
			return
		}

		utils.JSONSuccessResponse(c, "getting pre-register list", registerList)
	})

	r.POST("/register/:gameId", func(ctx *gin.Context) {
		gameId := ctx.Param("gameId")
		gameIdRegex, gameIdErr := regexp.Compile(`^\d*$`)
		if gameIdErr != nil {
			utils.JSONErrorResponse(ctx, 500, "inserting register gameId", gameIdErr)
			return
		}
		if !gameIdRegex.MatchString(gameId) {
			utils.JSONErrorResponse(ctx, 400, "inserting register gameId", errors.New("Unsupported parameters."))
			return
		}

		err := database.Register(gameId)
		if err != nil {
			utils.JSONErrorResponse(ctx, 500, "inserting register gameId", err)
			return
		}

		utils.JSONSuccessResponse(ctx, "inserting register gameId", nil)
	})

	// 包裝 Gin 以供 Lambda 使用
	ginLambda = ginadapter.NewV2(r)
}

// handling CORS.
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Next()
	}
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return ginLambda.ProxyWithContext(ctx, request)

	// if strings.Contains(request.Path, "/count/set") {
	// 	return handleSetCount(ctx, request)
	// } else if strings.Contains(request.Path, "/count/add") {
	// 	return handleAddCount(ctx, request)
	// } else if strings.Contains(request.Path, "/count") {
	// 	return handleGetCount(ctx, request)
	// } else if strings.Contains(request.Path, "/pre-register") {
	// 	return handlePreRegister(ctx, request)
	// } else {
	// 	return events.APIGatewayProxyResponse{
	// 		Headers:    headers,
	// 		StatusCode: 404,
	// 		Body:       "Not Found: " + request.Path,
	// 	}, nil
	// }
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
