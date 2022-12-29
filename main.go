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
	corsConfig.AllowOrigins = []string{"http://localhost:5173", "http://127.0.0.1:5173"}
	// corsConfig.AllowAllOrigins = true
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
	app.GET("/student", middlewares.RequireAuth, controllers.RetrieveStudentProfile)
	app.GET("/student/courses/:studID", middlewares.RequireAuth, controllers.RetrieveEnrolledCourses)
	app.GET("/student/courses/:studID/per", middlewares.RequireAuth, controllers.GetByDay)

	// Course
	app.GET("/course/:courseID", middlewares.RequireAuth, controllers.RetrieveCourse)

	// REGISTRAR
	app.GET("/admin", middlewares.RequireAuth, controllers.RetrieveAdminProfile)
	app.GET("/admin/courses", middlewares.RequireAuth, controllers.AllCourses)
	app.GET("/admin/students", middlewares.RequireAuth, controllers.RetrieveStudents)

	app.POST("/admin/courses/add", middlewares.RequireAuth, controllers.InsertCourse)
	app.POST("/admin/student/add", middlewares.RequireAuth, controllers.RegisterStudent)
	app.POST("/admin/student/enroll", middlewares.RequireAuth, controllers.EnrollStudent)

	app.DELETE("/admin/student/:studID", middlewares.RequireAuth, controllers.DeleteStudent)
	app.DELETE("/admin/course/:courseID", middlewares.RequireAuth, controllers.DeleteCourse)
	app.DELETE("/admin/drop/:studID/:courseID", middlewares.RequireAuth, controllers.StudentDropCourse)

	app.POST("/admin/student/enrolls", controllers.BulkEnrollStudent)

	// MASTER
	app.POST("/master/registrar/create", controllers.MasterCreateRegistrar)

	// TESTS
	app.GET("/test/students", tests.TestGetAllStudents)
	app.GET("/test/validate/students", middlewares.RequireAuth, tests.TestStudentValidate)
	app.GET("/test/validate/registrar", middlewares.RequireAuth, tests.TestRegisValidate)

	app.Run()
}
