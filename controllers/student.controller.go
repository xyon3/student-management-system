package controllers

import (
	"encoding/json"
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

	c.JSON(http.StatusOK, gin.H{
		"msg":    "login successfully",
		"status": 200,
		"token":  tokenString,
	})
}

func RetrieveStudentProfile(c *gin.Context) {

	if role, _ := c.Get("role"); role != "student" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	contextStudID, _ := c.Get("studID")

	var student models.Student

	if err := models.DB.Select("studID, name, profileImg, ringDelay").Where("studID = ?", contextStudID).First(&student).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"msg":    "student does not exist",
			"err":    err,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   student,
	})

}

// todo: sort functionality
func RetrieveEnrolledCourses(c *gin.Context) {
	// STUDENT AUTH
	if role, _ := c.Get("role"); role != "student" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	contextStudID, _ := c.Get("studID")

	if c.Param("studID") == contextStudID {

		var enrolled []models.StudentEnrolledCourses

		models.DB.Table("tblCourse").Select("tblCourse.*, tblStudent.ringDelay").Joins("JOIN juncEnrolled ON tblCourse.courseID = juncEnrolled.courseID").Joins("JOIN tblStudent ON tblStudent.studID = juncEnrolled.studID").Where("juncEnrolled.studID = ?", contextStudID).Find(&enrolled)

		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"data": enrolled,
		})
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"msg": "unauthorized token",
		})
	}
}

func GetByDay(c *gin.Context) {
	if role, _ := c.Get("role"); role != "student" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	contextStudID, _ := c.Get("studID")

	if c.Param("studID") == contextStudID {

		var enrolled []models.StudentEnrolledCourses

		models.DB.Table("tblCourse").Select("tblCourse.*, tblStudent.ringDelay").Joins("JOIN juncEnrolled ON tblCourse.courseID = juncEnrolled.courseID").Joins("JOIN tblStudent ON tblStudent.studID = juncEnrolled.studID").Where("juncEnrolled.studID = ?", contextStudID).Find(&enrolled)

		type studentEnrolled struct {
			CourseID    string   `json:"courseID" gorm:"primaryKey;column:courseID"`
			Description string   `json:"description"`
			Proctor     string   `json:"proctor"`
			Day         []string `json:"day"`
			StartTime   string   `json:"startTime" gorm:"column:startTime"`
			EndTime     string   `json:"endTime" gorm:"column:endTime"`
			RingDelay   string   `json:"ringDelay" gorm:"column:ringDelay"`
			Room        string   `json:"roomLoc" gorm:"column:roomLoc"`
		}

		var enrolledFixedDay []studentEnrolled

		for i := 0; i < len(enrolled); i++ {
			var fixedDay []string
			json.Unmarshal([]byte(enrolled[i].Day), &fixedDay)
			enrolledFixedDay = append(enrolledFixedDay, studentEnrolled{
				CourseID:    enrolled[i].CourseID,
				Description: enrolled[i].Description,
				Proctor:     enrolled[i].Proctor,
				Day:         fixedDay,
				StartTime:   enrolled[i].StartTime,
				EndTime:     enrolled[i].EndTime,
				RingDelay:   enrolled[i].RingDelay,
				Room:        enrolled[i].Room,
			})
		}

		// c.JSON(http.StatusOK, enrolledFixedDay)
		// fmt.Println(len(enrolledFixedDay))

		type perDay struct {
			Day     string            `json:"day"`
			Courses []studentEnrolled `json:"courses"`
		}

		var coursesPerDay []perDay

		var DAYS = []string{"M", "T", "W", "TH", "F", "S"}

		for currentDay := 0; currentDay < len(DAYS); currentDay++ {
			var dayCourses []studentEnrolled
			for currentCourse := 0; currentCourse < len(enrolledFixedDay); currentCourse++ {
				for currentDayCourse := 0; currentDayCourse < len(enrolledFixedDay[currentCourse].Day); currentDayCourse++ {
					if enrolledFixedDay[currentCourse].Day[currentDayCourse] == DAYS[currentDay] {
						dayCourses = append(dayCourses, enrolledFixedDay[currentCourse])
					}
				}
			}
			coursesPerDay = append(coursesPerDay, perDay{
				Day:     DAYS[currentDay],
				Courses: dayCourses,
			})

		}
		switch queryDay := c.Query("day"); queryDay {
		case "monday":
			c.JSON(http.StatusOK, coursesPerDay[0])
		case "teusday":
			c.JSON(http.StatusOK, coursesPerDay[1])
		case "wednesday":
			c.JSON(http.StatusOK, coursesPerDay[2])
		case "thursday":
			c.JSON(http.StatusOK, coursesPerDay[3])
		case "friday":
			c.JSON(http.StatusOK, coursesPerDay[4])
		case "saturday":
			c.JSON(http.StatusOK, coursesPerDay[5])
		case "all":
			c.JSON(http.StatusOK, coursesPerDay)
		default:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"msg":    "invalid params",
				"status": http.StatusBadRequest,
			})
		}

	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"msg": "unauthorized token",
		})
	}

}

// FUNTCION: updateSudentProfile() [this can be accessed by the registrar]
