package main

import (
	"course-api/internal/database"
	"course-api/internal/exercise/usecase"
	"course-api/internal/middlewares"
	userUc "course-api/internal/user/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	db := database.NewDatabaseConn()
	exerciseUsecase := usecase.NewExerciseUsecase(db)
	userUcs := userUc.NewUserUseCase(db)
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	r.GET("/exercises/:id", middlewares.WithAuth(userUcs), exerciseUsecase.GetExercise)
	r.GET("/exercises/:id/scores", middlewares.WithAuth(userUcs), exerciseUsecase.CalculateScore)
	r.POST("/exercises", middlewares.WithAuth(userUcs), exerciseUsecase.CreateExercise)
	r.POST("/exercises/questions", middlewares.WithAuth(userUcs), exerciseUsecase.CreateQuestion)
	r.POST("/exercises/:ExercisesId/questions/:QuestionsId/answer", middlewares.WithAuth(userUcs), exerciseUsecase.CreateAnswer)

	r.POST("/register", userUcs.Register)
	r.POST("/login", userUcs.Login)

	r.Run(":5000")
}
