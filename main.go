package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"man-bank/handlers"
	"man-bank/models"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Connect to the database using GORM
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", dbHost, dbUser, dbPassword, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	// Auto migrate models
	db.AutoMigrate(&models.User{}, &models.Account{}, &models.Transaction{})

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Validator
	e.Validator = handlers.NewValidator()

	// Routes
	authHandler := &handlers.AuthHandler{DB: db}
	e.POST("/signup", authHandler.Signup)
	e.POST("/login", authHandler.Login)

	// JWT Middleware config
	jwtConfig := middleware.JWTConfig{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}

	// Protected routes
	e.GET("/protected", func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		userId, ok := claims["user_id"].(float64)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid user ID"})
		}
		username, ok := claims["username"].(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid username"})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":  "You are in a protected route!",
			"user_id":  userId,
			"username": username,
		})
	}, middleware.JWTWithConfig(jwtConfig))

	depositHandler := &handlers.DepositHandler{DB: db}
	e.POST("/deposit", depositHandler.Deposit, middleware.JWTWithConfig(jwtConfig))

	accountHistoryHandler := &handlers.AccountHistoryHandler{DB: db}
	e.GET("/account/:account_id/history", accountHistoryHandler.GetAccountHistory, middleware.JWTWithConfig(jwtConfig))

	withdrawHandler := &handlers.WithdrawHandler{DB: db}
	e.POST("/withdraw", withdrawHandler.Withdraw, middleware.JWTWithConfig(jwtConfig))

	transferHandler := &handlers.TransferHandler{DB: db}
	e.POST("/transfer", transferHandler.Transfer, middleware.JWTWithConfig(jwtConfig))
	transferHistoryHandler := &handlers.TransferHistoryHandler{DB: db}
	e.GET("/account/:account_id/transfer-history", transferHistoryHandler.GetTransferHistory, middleware.JWTWithConfig(jwtConfig))
	e.Logger.Fatal(e.Start(":8080"))
}
