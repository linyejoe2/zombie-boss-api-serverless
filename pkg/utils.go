// Version 0.0.1 by linyejoe2 at 12/14/24 04:38

package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type IResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Body    interface{} `json:"body"`
	// Code    int         `json:"code,omitempty"`
}

// JSONResponse 用來回傳標準化的 JSON 回應格式。
//
// Parameters:
//   - ctx (*gin.Context): Gin 框架的 Context，負責處理 HTTP 請求及回應。
//   - statusCode (int): HTTP 回應狀態碼 (例如：200、400、500)。
//   - err (bool): 用來表示回應中是否有錯誤，true 代表有錯誤，false 代表正常回應。
//   - message (string): 回應的訊息，通常用來描述成功或失敗的狀況。
//   - body (interface{}): 回應的主要內容，可以是任何資料型別。
//
// Example:
//
//	JSONResponse(ctx, 200, false, "Do something successfully", data)
func JSONResponse(ctx *gin.Context, statusCode int, err bool, message string, body interface{}) {
	response := IResponse{
		Error:   err,
		Message: message,
		Body:    body,
		// Code:    statusCode,
	}
	ctx.JSON(statusCode, response)
}

// JSONErrorResponse 用來回傳操作失敗時的 JSON 錯誤訊息。
//
// Parameters:
//   - ctx (*gin.Context): Gin 框架的 Context，負責處理 HTTP 請求及回應。
//   - statusCode (int): HTTP 回應狀態碼 (例如：400、500)，代表失敗狀態。
//   - operation (string): 描述所嘗試的操作名稱 (例如："create user", "fetch data")。
//   - err (error): 具體的錯誤訊息，會被包含在回應中。
//
// Example:
//
//	JSONErrorResponse(ctx, 500, "fetch data", err)
func JSONErrorResponse(ctx *gin.Context, statusCode int, operation string, err error) {
	JSONResponse(ctx, statusCode, true, fmt.Sprintf("Failed to %v: %v", operation, err), nil)
}

// JSONSuccessResponse 用來回傳操作成功時的 JSON 資料。
//
// Parameters:
//   - ctx (*gin.Context): Gin 框架的 Context，負責處理 HTTP 請求及回應。
//   - operation (string): 描述所嘗試的操作名稱 (例如："create user", "fetch data")。
//   - body (interface{}): 回應的主要內容，可以是任何資料型別。
//
// Example:
//
//	JSONSuccessResponse(ctx, "fetch data", err)
func JSONSuccessResponse(ctx *gin.Context, operation string, body interface{}) {
	JSONResponse(ctx, 200, false, fmt.Sprintf("%v successfully.", operation), body)
}
