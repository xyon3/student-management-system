package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sesh/models"
)

func RetrieveCourse(c *gin.Context) {

	if role, _ := c.Get("role"); role != "admin" && role != "student" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var course models.Course
	models.DB.Where("courseID = ?", c.Param("courseID")).Find(&course)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   course,
	})
}
