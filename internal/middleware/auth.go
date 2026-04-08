package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pisira/backend/pkg/response"
)

func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		println(tokenStr)
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		// Simpan info user ke context agar bisa dipakai di handler
		c.Set("user_id", int(claims["user_id"].(float64)))
		c.Set("user_role", claims["role"].(string))
		c.Next()
	}
}

// AdminOnly memastikan hanya role admin yang bisa akses
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("user_role")
		if role != "admin" {
			c.JSON(403, gin.H{"success": false, "message": "Hanya admin yang dapat mengakses fitur ini"})
			c.Abort()
			return
		}
		c.Next()
	}
}
