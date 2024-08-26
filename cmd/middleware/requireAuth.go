package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Get cookie
		tokenString, err := c.Cookie("Authorization")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte("secret-key"), nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		//parse claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if float64(time.Now().Unix()) >= claims["exp"].(float64) { // проверка exp токена
				c.AbortWithStatus(http.StatusUnauthorized) // exp время токена конец
			}
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Next()
	}
}
