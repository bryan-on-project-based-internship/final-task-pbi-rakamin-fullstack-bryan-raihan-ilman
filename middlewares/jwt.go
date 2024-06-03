package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"user-profile-apis/helpers"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		fmt.Println("Received Authorization header:", authHeader)
		if authHeader == "" {
			fmt.Println("No Authorization header provided")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := strings.SplitN(authHeader, " ", 2)[1]
		fmt.Println("Attempting to validate JWT...")
		token, err := helpers.ValidateJWT(tokenString)
		if err != nil {
			fmt.Println("Error validating JWT:", err)
			if err == jwt.ErrSignatureInvalid {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if !token.Valid {
			fmt.Println("Invalid JWT")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := claims["uid"].(float64)
		c.Set("userID", uint(userID)) // Set user ID for access in handlers

		fmt.Println("JWT validated successfully. Proceeding to next handler...")
		c.Next()
	}
}