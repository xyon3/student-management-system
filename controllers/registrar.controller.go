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

func RegistrarLogin(c *gin.Context) {

	var reqBody models.RegistrarLoginBody
	// BIND REQUEST TO BODY
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "something went wrong. could not bind request body",
			"err": err.Error(),
		})
		return
	}

	// QUERY USER
	var dbData models.Registrar
	if res := models.DB.Select("regID", "hash").Where("regID = ?", reqBody.RegID).First(&dbData); res.Error != nil {
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
		"iss": dbData.RegID,
		"sub": "reg",
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

	c.JSON(http.StatusOK, gin.H{
		"msg":    "login successfully",
		"status": 200,
		"token":  tokenString,
	})
}

func RegisterStudent(c *gin.Context) {

	// CHECK IF AUTH IS ADMIN
	if role, _ := c.Get("role"); role != "admin" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"": "",
		})
		return
	}

	// ATTACH REQUEST BODY TO MODEL
	var reqBody models.Student
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": "could not attach to request body",
		})
		return
	}

	// CHECK IF STUDENT ALREADY EXIST
	var dbData models.Student
	models.DB.Where("studID = ?", reqBody.StudID).First(&dbData)

	if dbData.StudID != "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": "entity already exist",
		})
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(reqBody.Hash), bcrypt.DefaultCost)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": "could not hash the password",
		})
		return
	}

	reqBody.Hash = string(hashedPass)

	// INSERT REQ BODY DATA TO tblStudent TABLE
	if result := models.DB.Create(&reqBody); result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": "DB failed",
		})
		return
	}

	// DONE
	c.AbortWithStatusJSON(http.StatusCreated, gin.H{
		"msg": "entity created",
	})
}

func EnrollStudent(c *gin.Context) {
	if role, _ := c.Get("role"); role != "admin" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// DO MORE

	// INIT REQUEST BODY
	var reqBody models.EnrolledRequestBody

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": "could not bind request body",
		})
		return
	}

	var dbData models.Enrolled
	// CHECK IF STUDENT IS ALREADY ENROLLED IN THE COURSE
	if result := models.DB.Where("studID = ? AND courseID = ?", reqBody.StudID, reqBody.CourseID).Find(&dbData); result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": "entity already exist",
		})
		return
	}

	// ENROLL STUDENT
	if result := models.DB.Create(&reqBody); result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": "could not create entity",
		})
		return
	}

	// SUCCESS
	c.AbortWithStatusJSON(http.StatusCreated, gin.H{
		"msg": "entity created",
	})
}

func BulkEnrollStudent(c *gin.Context) {
	// if role, _ := c.Get("role"); role != "admin" {
	// 	c.AbortWithStatus(http.StatusUnauthorized)
	// 	return
	// }

	// DO MORE

	// INIT REQUEST BODY
	var reqBody []models.EnrolledRequestBody

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": "could not bind request body",
		})
		return
	}
	// c.JSON(200, gin.H{
	// 	"data": reqBody,
	// })
	// fmt.Println(len(reqBody))

	var dbData models.Enrolled
	for i := 1; i <= len(reqBody); i++ {

		c.JSON(200,
			reqBody[i],
		)

		if result := models.DB.Where("studID = ? AND courseID = ?", reqBody[i].StudID, reqBody[i].CourseID).Find(&dbData); result.Error != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"msg": "entity already exist",
			})
			return
		}
		if result := models.DB.Create(&reqBody[i]); result.Error != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"msg": "could not create entity",
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"msg":      "entity created",
			"affected": i,
		})
	}

}

func InsertCourse(c *gin.Context) {
	// CHECK IF AUTH IS ADMIN

	if role, _ := c.Get("role"); role != "admin" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// ATTACH REQUEST BODY TO MODEL
	var reqBody models.Course

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		// c.AbortWithStatus(http.StatusInternalServerError)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": "could not bind to reqBody",
		})
		return
	}

	// CHECK IF COURSE ALREADY EXIST
	var dbData models.Course

	if result := models.DB.Where("courseID = ?", reqBody.CourseID).First(&dbData); result.Error == nil {
		// c.AbortWithStatus(http.StatusInternalServerError)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": "course already exist",
		})
		return
	}
	// INSERT REQ BODY DATA TO tblCourse TABLE
	if err := models.DB.Create(&reqBody); err.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": "could not insert course",
			"err": err.Error,
		})
		return
	}
	// DONE
	c.JSON(http.StatusCreated, gin.H{
		"msg": "entity created",
	})
}
