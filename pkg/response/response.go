package response

import "github.com/gin-gonic/gin"

type ResponseInterface interface {
	SendError(c *gin.Context, statusCode int, message string)
	Success(c *gin.Context, statusCode int, data interface{})
}

type Response struct{}

type ErrorResponse struct {
	Message string `json:"message"`
	Error   bool   `json:"error"`
}

type SuccessResponse struct {
	Data  interface{} `json:"data"`
	Error bool        `json:"error"`
}

func NewResponse() ResponseInterface {
	return &Response{}
}

func (r Response) SendError(c *gin.Context, statusCode int, message string) {
	response := ErrorResponse{
		Error:   true,
		Message: message,
	}
	c.AbortWithStatusJSON(statusCode, response)
}

func (r Response) Success(c *gin.Context, statusCode int, data interface{}) {
	response := SuccessResponse{
		Error: false,
		Data:  data,
	}
	c.JSON(statusCode, response)
}
