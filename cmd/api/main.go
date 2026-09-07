package main

import (
	"log"
	"time"

	"medix-be/config"
	"medix-be/migrations"

	"medix-be/internal/drug"
	"medix-be/internal/medicine"
	"medix-be/internal/middleware"
	"medix-be/internal/report"
	"medix-be/internal/transaction"
	"medix-be/internal/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	config.LoadEnv()
	migrations.RunMigrations()
	config.ConnectDatabase()

	logger := logrus.New()
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // URL Frontend Vite kamu
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	apiV1 := r.Group("/api/v1")

	authAPI := apiV1.Group("")
	authAPI.Use(middleware.Auth())

	medicine.StartApp(&medicine.ModuleConfig{
		DB:     config.DB,
		Logger: logger,
		Router: authAPI,
	})

	drug.StartApp(&drug.TypeDrugHandler{
		DB:     config.DB,
		Logger: logger,
		Router: authAPI,
	})

	user.StartApp(&user.UserHandler{
		DB:     config.DB,
		Logger: logger,
		Router: apiV1,
		AuthRouter:   authAPI,
	})

	transaction.StartApp(&transaction.TransactionHandler{
		DB:     config.DB,
		Logger: logger,
		Router: authAPI,
	})

	report.StartApp(&report.ReportHandler{
		DB:     config.DB,
		Logger: logger,
		Router: authAPI,
	})

	log.Println("Server Medix BE running on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
