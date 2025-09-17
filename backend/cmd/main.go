package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"account-management/backend/internal/repository"
	"account-management/backend/pkg/config"
)

func main() {
	// 環境変数を読み込み
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// DB接続
	dbCfg := config.LoadDB()
	db, err := repository.InitDB(dbCfg.DSN())
	if err != nil {
		log.Fatalf("failed to connect DB %v", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	// Ginルーターを設定
	r := gin.Default()

	// CORS設定
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// ヘルスチェック（アプリ / DB）
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Account Management API is running",
		})
	})

	r.GET("/health/db", func(c *gin.Context){
		if err := db.Ping(); err != nil {
			c.JSON(500, gin.H{"db": "down", "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"db": "up"})
	})



	// 基本的なAPIエンドポイント
	api := r.Group("/api/v1")
	{
		api.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "API is working!",
			})
		})
	}

	// デバッグ用：起動時に全ルートをログ出力
    for _, rt := range r.Routes() {
        log.Printf("[ROUTE] %s %s", rt.Method, rt.Path)
    }

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}