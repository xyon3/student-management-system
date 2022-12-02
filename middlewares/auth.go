package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sesh/models"
)

func RequireAuth(c *gin.Context) {

	// get cookie from Authorization request header
	authHeader := c.Request.Header["Authorization"]

	if authHeader == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tokenString := authHeader[0]

	if tokenString == "" {
		c.AbortWithStatus(http.StatusUnauthorized)

		// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		// 	"msg": "token not found",
		// })
		return
	}

	// Decode and validate the supplied jwt
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		// c.AbortWithStatus(http.StatusInternalServerError)

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": "could not parse token",
		})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		if claims["sub"] == "stud" {
			// find the user with token sub
			var student models.Student
			models.DB.Select("studID", "name", "profileImg", "ringDelay").Where("studID = ?", claims["iss"]).Find(&student)

			if student.StudID == "" {
				c.AbortWithStatus(http.StatusUnauthorized)
			}
			// Attach identity to request
			c.Set("student", student)
			c.Set("role", "student")
		} else if claims["sub"] == "reg" {

			var registrar models.Registrar
			models.DB.Select("regID", "name").Where("regID = ?", claims["iss"]).Find(&registrar)

			// fmt.Println(registrar)

			if registrar.RegID == "" {
				c.AbortWithStatus(http.StatusUnauthorized)
			}
			c.Set("registrar", registrar)
			c.Set("role", "admin")
		}

		// CONTINUE
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
