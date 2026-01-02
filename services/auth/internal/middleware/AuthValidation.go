package middleware

import (
	"net/http"
	"time"

	"github.com/adamhaiqal/go-auth/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthenticateCookie(c *gin.Context) {
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return utils.GetJWTKey(), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		if exp, ok := claims["exp"].(float64); ok {
			if float64(time.Now().Unix()) > exp {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var accountID uint
		switch v := claims["sub"].(type) {
		case float64:
			accountID = uint(v)
		case int:
			accountID = uint(v)
		case uint:
			accountID = v
		default:
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("accountID", accountID)
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
	}
}
