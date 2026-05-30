package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mahmoud-abd-elnasser/todo_api/internal/config"
	"github.com/mahmoud-abd-elnasser/todo_api/internal/database"
	"github.com/mahmoud-abd-elnasser/todo_api/internal/handlers"
	"github.com/mahmoud-abd-elnasser/todo_api/internal/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	pool, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}

	defer pool.Close()

	router := gin.Default()

	router.SetTrustedProxies(nil)

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message":  "ToDo API is running!",
			"status":   "success",
			"database": "connected",
		})
	})

	router.POST("/auth/register", handlers.CreateUserHandler(pool))
	router.POST("/auth/login", handlers.LoginHandler(pool, cfg))

	protectedUsers := router.Group("/users")
	protectedUsers.Use(middleware.AuthMiddleware(cfg))

	{
		protectedUsers.GET("", handlers.GetUserByEmailHandler(pool))
		protectedUsers.GET("/:id", handlers.GetUserByIdHandler(pool))
	}

	protectedTodos := router.Group("/todos")
	protectedTodos.Use(middleware.AuthMiddleware(cfg))

	{
		protectedTodos.GET("", handlers.GetAllTodosHandler(pool))
		protectedTodos.POST("", handlers.CreateTodoHandler(pool))
		protectedTodos.GET("/:id", handlers.GetTodoByIdHandler(pool))
		protectedTodos.PUT("/:id", handlers.UpdateTodoHandler(pool))
		protectedTodos.DELETE("/:id", handlers.DeleteTodoHandler(pool))
	}

	router.GET("/auth/protected", middleware.AuthMiddleware(cfg), handlers.TestProtectedHandler())

	router.Run(":" + cfg.Port)
}
