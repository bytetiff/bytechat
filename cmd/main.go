package main

import (
	"log"
	"os"

	"github.com/bytetiff/bytechat/internal/db"
	"github.com/bytetiff/bytechat/internal/handlers"
	"github.com/bytetiff/bytechat/internal/middleware"
	"github.com/bytetiff/bytechat/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Warning: .env file not found, using system environment variables")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET not set")
	}
	pool, err := db.ConnectPostgres(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Создаём репозитории
	userRepo := repository.NewUserRepository(pool)
	chatRepo := repository.NewChatRepository(pool)
	msgRepo := repository.NewMessageRepository(pool)
	refreshRepo := repository.NewRefreshTokenRepository(pool)

	// Создаём хендлеры
	authHandler := handlers.NewAuthHandler(userRepo, refreshRepo, []byte(jwtSecret))
	chatHandler := handlers.NewChatHandler(chatRepo)
	messageHandler := handlers.NewMessageHandler(msgRepo, chatRepo)

	r := gin.Default()

	// Публичные маршруты
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.POST("/refresh", authHandler.RefreshToken)
	r.POST("/logout", authHandler.Logout)

	// Middleware
	authMiddleware := middleware.AuthMiddleware(refreshRepo, []byte(jwtSecret))

	// Защищённые роуты
	protected := r.Group("/")
	protected.Use(authMiddleware)
	{
		protected.GET("/profile", authHandler.Profile)

		// Чаты
		protected.POST("/chats", chatHandler.CreateChat)
		protected.GET("/chats", chatHandler.GetUserChats)
		protected.POST("/chats/:chat_id/members", chatHandler.AddMember)

		// Сообщения
		protected.POST("/chats/:chat_id/messages", messageHandler.CreateMessage)
		protected.GET("/chats/:chat_id/messages", messageHandler.GetMessages)
	}

	log.Println("ByteChat server is running on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
