package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mahmoud-abd-elnasser/todo_api/internal/config"
	"github.com/mahmoud-abd-elnasser/todo_api/internal/models"
	"github.com/mahmoud-abd-elnasser/todo_api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string `json:"token"`
}

func CreateUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input RegisterRequest

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if input.Email == "" || input.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
			return
		}

		if len(input.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters long"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user , err := repository.CreateUser(pool, &models.User{
			Email: input.Email,
			Password: string(hashedPassword),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, user)
	}
}


func GetUserByIdHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		
		strId := c.Param("id")

		// id, err := strconv.Atoi(strId)
		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		// 	return
		// }
		if strId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
			return
		}

		user , err := repository.GetUserById(pool, strId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, user)
	}
}


func GetUserByEmailHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		
		email := c.Query("email")

		user , err := repository.GetUserByEmail(pool, email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func LoginHandler(pool *pgxpool.Pool , cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest LoginRequest

		if err := c.BindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		user, err := repository.GetUserByEmail(pool, loginRequest.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid email or password",
			})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid email or password",
			})
			return
		}

		claims := jwt.MapClaims{
			"user_id": user.ID,
			"email": user.Email,
			"exp":  time.Now().Add(24 * time.Hour).Unix(),

		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString([]byte(cfg.JWTSecret))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate token",
			})
			return
		}

		c.JSON(http.StatusOK, LoginResponse{
			Token: tokenString,
		})
	}
}

func TestProtectedHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId , exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError , gin.H{
				"error": "user_id not found in context",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, user with ID: " + userId.(string),
			"user_id": userId,
		})
	}
}
