package middleware

import (
	"RestAPI/internal/user"
	"context"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	bearer = "Bearer "
)

//type contextKey string

func Authentication(userService *user.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{
				"message:": "sorry not authorize",
			})
			c.Abort()
		}
		if !strings.HasPrefix(authHeader, bearer) {
			c.JSON(401, gin.H{
				"message:": "sorry not authorize",
			})
			c.Abort()
		}
		auths := strings.Split(authHeader, " ") //will result []string{"Bearer", "{token}"}
		data, err := userService.DecriptJWT(auths[1])
		//fmt.printf("%+v\n", data)
		if err != nil {
			c.JSON(401, gin.H{
				"message:": "sorry not authorize",
			})
			c.Abort()
		}
		//context.Context	mendapatkan user id siapa dari value mana
		ctxUserID := context.WithValue(c.Request.Context(), "user_id", data["user_id"])
		c.Request = c.Request.WithContext(ctxUserID)
		c.Next()
	}
}
