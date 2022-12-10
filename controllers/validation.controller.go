package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidationResult(c *gin.Context) {
	role, _ := c.Get("role")

	if role != "admin" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"role": role,
	})
}
