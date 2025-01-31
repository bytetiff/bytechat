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
	// Загружаем .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Warning: .env file not found, using system environment variables")
	}

	// Читаем переменные из окружения
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	jwtSecret := os.Getenv("JWT_SECRET")

	// Проверяем, что важные переменные загружены
	if jwtSecret == "" {
		log.Fatal("❌ ERROR: JWT_SECRET is required but not set")
	}

	// Подключаемся к базе данных
	pool, err := db.ConnectPostgres(dbHost, dbPort, dbUser, dbPass, dbName)
	if err != nil {
		log.Fatal("❌ ERROR: failed to connect to database:", err)
	}
	defer pool.Close()

	// Создаём репозиторий пользователей
	userRepo := repository.NewUserRepository(pool)

	// Создаём хендлеры с JWT
	refreshTokenRepo := repository.NewRefreshTokenRepository(pool)
	authHandler := handlers.NewAuthHandler(userRepo, refreshTokenRepo, []byte(jwtSecret))
	// Подключение middleware для аутентификации
	authMiddleware := middleware.AuthMiddleware(refreshTokenRepo, []byte(jwtSecret))

	// Группа защищённых маршрутов
	r := gin.Default()
	protected := r.Group("/")
	protected.Use(authMiddleware)
	protected.GET("/profile", authHandler.Profile)

	// Запускаем сервер
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.POST("/refresh", authHandler.RefreshToken)
	r.POST("/logout", authHandler.Logout)

	log.Println("🚀 ByteChat server is running on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
