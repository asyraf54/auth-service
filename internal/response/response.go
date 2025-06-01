package response

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ErrorDetail struct {
	ErrorCode         string `json:"errorCode"`
	ErrorMessage      string `json:"errorMessage"`
	ErrorDebugMessage string `json:"errorDebugMessage"`
}

type APIResponse struct {
	TransactionID string      `json:"transactionId"`
	StatusCode    int         `json:"statusCode"`
	Data          interface{} `json:"data,omitempty"`
	Error         *ErrorDetail `json:"error,omitempty"`
}

func JSON(c *gin.Context, status int, data interface{}, errDetail *ErrorDetail) {
	res := APIResponse{
		TransactionID: uuid.NewString(),
		StatusCode:    status,
		Data:          data,
		Error:         errDetail,
	}
	c.JSON(status, res)
}

func AbortWithStatusJSON(c *gin.Context, status int, data interface{}, errDetail *ErrorDetail) {
	res := APIResponse{
		TransactionID: uuid.NewString(),
		StatusCode:    status,
		Data:          data,
		Error:         errDetail,
	}
	c.AbortWithStatusJSON(status, res)
}
