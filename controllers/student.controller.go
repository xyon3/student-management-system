package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sesh/models"
	"golang.org/x/crypto/bcrypt"
)

func StudentLogin(c *gin.Context) {

	var reqBody models.StudentLoginBody
	// BIND REQUEST TO BODY
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "something went wrong. could not bind request body",
			"err": err.Error(),
		})
		return
	}

	// QUERY USER
	var dbData models.StudentLoginBody
	if res := models.DB.Select("studID", "hash").Where("studID = ?", reqBody.StudID).First(&dbData); res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "user does not exist",
		})
		return
	}

	// COMPARE HASH
	if err := bcrypt.CompareHashAndPassword([]byte(dbData.Hash), []byte(reqBody.Hash)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "Invalid data supplied",
		})
		return
	}

	// GENERATE JWT TOKEN

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": dbData.StudID,
		"sub": "stud",
		"exp": time.Now().Add(time.Hour * 6).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "failed to create token",
			"err": err.Error(),
		})
		return
	}
	//////// DONE
	// SEND IT BACK VIA CROSS ORIGIN REQUESTS
	// DONE: MAKE THE SERVER STATELESS (REMOVE COOKIES ON SERVER WITHOUT BREAKING AUTH)
	// JWT AS SESSION IS BAD *pouting emoji*

	// c.SetSameSite(http.SameSiteLaxMode)
	// c.SetCookie("auth", tokenString, 3600*24, "", "", false, true)
	//////// DONE

	c.JSON(http.StatusOK, gin.H{
		"msg":    "login successfully",
		"status": 200,
		"token":  tokenString,
	})
}

func RetrieveEnrolledCourses(c *gin.Context) {
	// STUDENT AUTH
	if role, _ := c.Get("role"); role != "student" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	contextStudID, _ := c.Get("studID")

	var enrolled []models.StudentEnrolledCourses
	// RETRIVE ALL DATA WHERE studID == REQUEST BODY STUDENT ID
	models.DB.Table("tblCourse").Select("tblCourse.*, tblStudent.ringDelay").Joins("JOIN juncEnrolled ON tblCourse.courseID = juncEnrolled.courseID").Joins("JOIN tblStudent ON tblStudent.studID = juncEnrolled.studID").Where("juncEnrolled.studID = ?", contextStudID).Find(&enrolled)

	c.JSON(http.StatusOK, gin.H{
		"data": enrolled,
	})

	// SELECT * FROM tblCourses:w

	// INNER JOIN juncEnrolled
	// ON tblCourses.courseID = juncEnrolled.courseID
	// INNER JOIN tblStudent
	// ON tblStudent.studID = juncEnrolled.studID
	// DONE
}
