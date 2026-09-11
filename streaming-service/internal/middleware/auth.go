package middleware

import (
	"net/http"
	"stream-mesh/streaming/internal/payload"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			payload.Fail(c, http.StatusUnauthorized, "UNAUTHORIZED", "user unauthorized")
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		}, jwt.WithValidMethods([]string{"HS256"}))

		if err != nil || !token.Valid {
			payload.Fail(c, http.StatusUnauthorized, "UNAUTHORIZED", "user unauthorized")
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			payload.Fail(c, http.StatusUnauthorized, "UNAUTHORIZED", "user unauthorized")
			c.Abort()
			return
		}
		c.Set("user_id", claims["sub"])
		c.Next()
	}
}
