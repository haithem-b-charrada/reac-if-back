package middlewares

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID    int64
	Email string
}

func ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := User{}
		user.ID, _ = strconv.ParseInt(c.Request.Header.Get("User-Id"), 10, 64)
		user.Email = c.Request.Header.Get("User-Email")

		if user.ID <= 0 || user.Email == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "user not found",
				"status":  http.StatusUnauthorized,
			})

			return
		}

		c.Next()
	}
}
