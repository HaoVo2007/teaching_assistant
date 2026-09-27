package http

import (
	"teaching_assistant/internal/delivery/http/handler"
	"teaching_assistant/internal/delivery/http/middleware"
	userImport "teaching_assistant/internal/domain/user"
	"teaching_assistant/pkg/jwt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewRouter(
	app *fiber.App,
	userH *handler.UserHandler,
	questionH *handler.QuestionHandler,
	questionSetH *handler.QuestionSetHandler,
	classH *handler.ClassHandler,
	homeworkH *handler.HomeworkHandler,
	homeworkSubmissionH *handler.HomeworkSubmissionHandler,
	studentH *handler.StudentHandler,
	jwtManager *jwt.Manager,
) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173,https://teachingassistantfe.netlify.app",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	api := app.Group("/api/v1")
	admin := api.Group("/admin")
	admin.Use(middleware.RequireRole(jwtManager, userImport.RoleAdmin))
	// =========================auth routes=========================
	auth := api.Group("/auth")
	{
		auth.Post("/register", userH.Register)
		auth.Post("/login", userH.Login)
		auth.Post("/logout", middleware.AuthMiddleware(jwtManager), userH.Logout)
	}

	user := admin.Group("/users")
	{
		user.Post("/create", userH.CreateUser)
	}

	users := api.Group("/users")
	{
		users.Get("/parents", middleware.RequireRole(jwtManager, userImport.RoleTeacher), userH.GetParents)
	}
	// =========================auth routes=========================

	// =========================student routes=========================
	student := api.Group("/students")
	{
		student.Post("/claim", middleware.RequireRole(jwtManager, userImport.RoleTeacher), studentH.ClaimStudent)
		student.Get("/by-guardian", middleware.RequireRole(jwtManager, userImport.RoleParent), studentH.GetStudentByGuardian)
	}
	// =========================student routes=========================

	// =========================question routes=========================
	question := api.Group("/questions")
	question.Use(middleware.RequireRole(jwtManager, userImport.RoleTeacher))
	{
		question.Post("", questionH.CreateQuestion)
		question.Get("", questionH.GetQuestions)
		question.Get("/:id", questionH.GetQuestionById)
		question.Put("/:id", questionH.UpdateQuestionById)
		question.Delete("/:id", questionH.DeleteQuestionById)
	}
	// =========================question routes=========================

	// =========================question set routes=========================
	questionSet := api.Group("/question-sets")
	questionSet.Use(middleware.RequireRole(jwtManager, userImport.RoleTeacher))
	{
		questionSet.Post("", questionSetH.CreateQuestionSet)
		questionSet.Get("", questionSetH.GetQuestionSets)
		questionSet.Get("/:id", questionSetH.GetQuestionSetById)
		questionSet.Put("/:id", questionSetH.UpdateQuestionSetById)
		questionSet.Delete("/:id", questionSetH.DeleteQuestionSetById)
	}
	// =========================question set routes=========================

	// =========================class routes=========================
	class := api.Group("/classes")
	class.Use(middleware.RequireRole(jwtManager, userImport.RoleTeacher))
	{
		class.Post("", classH.CreateClass)
		class.Get("", classH.GetClasses)
		class.Get("/:id", classH.GetClassById)
		class.Put("/:id", classH.UpdateClassById)
		class.Delete("/:id", classH.DeleteClassById)
	}
	// =========================class routes=========================

	// =========================homework routes=========================
	homework := api.Group("/homeworks")
	{
		homework.Post("", middleware.RequireRole(jwtManager, userImport.RoleTeacher), homeworkH.CreateHomework)
		homework.Get("", middleware.RequireRole(jwtManager, userImport.RoleTeacher), homeworkH.GetHomeworks)
		homework.Get("/class/:class_id", middleware.RequireRole(jwtManager, userImport.RoleTeacher), homeworkH.GetHomeworksByClassId)
		homework.Get("/student/by-guardian", middleware.RequireRole(jwtManager, userImport.RoleParent), homeworkH.GetHomeworksByStudentId)
		homework.Get("/:id", middleware.RequireRole(jwtManager, userImport.RoleTeacher), homeworkH.GetHomeworkById)
		homework.Put("/:id", middleware.RequireRole(jwtManager, userImport.RoleTeacher), homeworkH.UpdateHomeworkById)
		homework.Delete("/:id", middleware.RequireRole(jwtManager, userImport.RoleTeacher), homeworkH.DeleteHomeworkById)
	}
	// =========================homework routes=========================

	// =========================homework submission routes=========================
	homeworkSubmission := api.Group("/homework-submissions")
	{
		homeworkSubmission.Post("", middleware.RequireRole(jwtManager, userImport.RoleParent), homeworkSubmissionH.CreateHomeworkSubmission)
		homeworkSubmission.Get("", middleware.RequireRole(jwtManager, userImport.RoleTeacher), homeworkSubmissionH.GetHomeworkSubmissions)
		homeworkSubmission.Get("/homework/:homework_id", middleware.RequireRole(jwtManager, userImport.RoleTeacher), homeworkSubmissionH.GetHomeworkSubmissionsByHomeworkId)
		homeworkSubmission.Get("/student/by-guardian/:homework_id", middleware.RequireRole(jwtManager, userImport.RoleParent), homeworkSubmissionH.GetHomeworkSubmissionsByHomeworkIdByGuardian)
		homeworkSubmission.Get("/:id", middleware.RequireRole(jwtManager, userImport.RoleTeacher), homeworkSubmissionH.GetHomeworkSubmissionById)
		// homeworkSubmission.Put("/:id", middleware.AuthMiddleware(jwtManager), homeworkSubmissionH.UpdateHomeworkSubmissionById)
		// homeworkSubmission.Delete("/:id", middleware.AuthMiddleware(jwtManager), homeworkSubmissionH.DeleteHomeworkSubmissionById)
	}
	// =========================homework submission routes=========================
}
