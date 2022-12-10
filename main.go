package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sesh/controllers"
	"github.com/sesh/middlewares"
	"github.com/sesh/models"
	"github.com/sesh/tests"
)

func main() {
	godotenv.Load()

	app := gin.Default()

	models.ConnectDB()

	corsConfig := cors.DefaultConfig()
	// corsConfig.AllowOrigins = []string{"http://localhost:5173", "http://127.0.0.1:5173"}
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Authorization"}
	corsConfig.AllowCredentials = true
	corsConfig.MaxAge = 3600 * 4

	app.Use(cors.New(corsConfig))

	// ROUTES:

	// AUTH
	app.POST("/auth/student", controllers.StudentLogin)
	app.POST("/auth/registrar", controllers.RegistrarLogin)
	app.POST("/auth/registrar/validate", middlewares.RequireAuth, controllers.ValidationResult)

	// STUDENT
	app.GET("/student/courses/enrolled", middlewares.RequireAuth, controllers.RetrieveEnrolledCourses)

	// REGISTRAR
	app.POST("/admin/courses/add", middlewares.RequireAuth, controllers.InsertCourse)
	app.POST("/admin/student/add", middlewares.RequireAuth, controllers.RegisterStudent)
	app.POST("/admin/student/enroll", middlewares.RequireAuth, controllers.EnrollStudent)

	app.POST("/admin/student/enrolls", controllers.BulkEnrollStudent)

	// MASTER
	app.POST("/master/registrar/create", controllers.MasterCreateRegistrar)

	// TESTS
	app.GET("/testAllStudents", tests.TestGetAllStudents)
	app.GET("/testStudentValidate", middlewares.RequireAuth, tests.TestStudentValidate)
	app.GET("/testRegisValidate", middlewares.RequireAuth, tests.TestRegisValidate)

	app.Run()
}
