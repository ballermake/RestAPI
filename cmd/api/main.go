package main

import (
	"RestAPI/internal/database"
	"RestAPI/internal/exercise"
	"RestAPI/internal/middleware"
	"RestAPI/internal/user"

	"github.com/gin-gonic/gin"
)

func main() {

	route := gin.Default()

	route.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]interface{}{
			"pesan": "Hello World",
		})
	})

	db := database.NewDatabaseConn()
	exerciseService := exercise.NewExerciseService(db)
	userService := user.NewUserService(db)
	//exercises
	//route.Use(middleware.Authentication(userService))		kalo begini dipakai buat general
	route.GET("/exercises/:id", middleware.Authentication(userService), exerciseService.GetExercise)
	route.GET("/exercises/:id/score", middleware.Authentication(userService), exerciseService.GetUserScore)

	//tugas
	route.POST("/exercises", middleware.Authentication(userService), exerciseService.CreateExercise)
	route.POST("/exercises/:id/question", middleware.Authentication(userService), exerciseService.CreateQuestion)
	route.POST("/exercises/:id/question/:qid/answer", middleware.Authentication(userService), exerciseService.CreateAnswer)

	//user
	route.POST("/register", userService.Register)
	route.POST("/login", userService.Login)

	route.Run(":8000")

}
