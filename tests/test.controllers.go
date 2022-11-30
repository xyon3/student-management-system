package tests

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sesh/models"
)

func TestStudentValidate(c *gin.Context) {
	user, _ := c.Get("student")
	c.JSON(http.StatusOK, gin.H{
		"msg": user,
	})
}

func TestRegisValidate(c *gin.Context) {
	user, _ := c.Get("registrar")
	c.JSON(http.StatusOK, gin.H{
		"msg": user,
	})
}

func TestGetAllStudents(c *gin.Context) {
	var students []models.Student
	models.DB.Find(&students)

	c.JSON(http.StatusOK, gin.H{
		"students": students,
		"nested": gin.H{
			"json": "inside a json",
		},
	})
}
