package middlewares

import (
	"net/http"
	"strings"

	"user-profile-apis/helpers"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := strings.SplitN(authHeader, " ", 2)[1]
		token, err := helpers.ValidateJWT(tokenString)
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := claims["uid"].(float64)
		c.Set("userID", uint(userID)) // Set user ID for access in handlers

		c.Next()
	}
}
