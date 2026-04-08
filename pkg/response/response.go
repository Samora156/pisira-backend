package response

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func OK(c *gin.Context, message string, data interface{}) {
	c.JSON(200, APIResponse{Success: true, Message: message, Data: data})
}

func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(201, APIResponse{Success: true, Message: message, Data: data})
}

func BadRequest(c *gin.Context, message string) {
	c.JSON(400, APIResponse{Success: false, Message: message})
}

func Unauthorized(c *gin.Context) {
	c.JSON(401, APIResponse{Success: false, Message: "Akses ditolak, silakan login terlebih dahulu"})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(404, APIResponse{Success: false, Message: message})
}

func ServerError(c *gin.Context, err error) {
	c.JSON(500, APIResponse{Success: false, Message: "Terjadi kesalahan server: " + err.Error()})
}
