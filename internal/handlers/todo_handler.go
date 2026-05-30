package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mahmoud-abd-elnasser/todo_api/internal/repository"
)

type CreateTodoInput struct {
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
}
type UpdateTodoInput struct {
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed" binding:"required"`
}

func CreateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}

		userID, ok := userIDInterface.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID in context is not a string"})
			return
		}

		var input CreateTodoInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		todo, err := repository.CreateTodo(pool, input.Title, input.Completed, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, todo)
	}
}

func UpdateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
				userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}

		userID, ok := userIDInterface.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID in context is not a string"})
			return
		}


		var input UpdateTodoInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError , gin.H{
			"error": err.Error(),
		})
		return
	}
		todo, err := repository.UpdateTodo(pool, id, input.Title, input.Completed, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, todo)

	}
}

func GetTodoByIdHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {

		id , err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format. Must be an integer."})
			return
		}

		todo, err := repository.GetTodoById(pool, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, todo)

	}
}

func GetAllTodosHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		todos, err := repository.GetAllTodos(pool)
		if err != nil {
			c.JSON(http.StatusInternalServerError , gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, todos)
	}
}

var ErrTodoNotFound = errors.New("todo not found or unauthorized")

func DeleteTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid ID. Please enter a valid integer ID",
			})
			return
		}

		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}

		userID, ok := userIDInterface.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal configuration error"})
			return
		}

		err = repository.DeleteTodo(pool, id, userID)
		if err != nil {
			if errors.Is(err, ErrTodoNotFound) {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Todo not found or you do not have permission to delete it",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to delete todo due to a server error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Todo deleted successfully",
		})
	}
}
