package main

import (
	"log"
	"time"

	"medix-be/config"
	"medix-be/migrations"

	// Import modul-modul kamu
	"medix-be/internal/drug"
	"medix-be/internal/medicine"
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

	medicine.StartApp(&medicine.ModuleConfig{
		DB:     config.DB,
		Logger: logger,
		Router: apiV1,
	})

	drug.StartApp(&drug.TypeDrugHandler{
		DB:     config.DB,
		Logger: logger,
		Router: apiV1,
	})

	user.StartApp(&user.UserHandler{
		DB:     config.DB,
		Logger: logger,
		Router: apiV1,
	})

	transaction.StartApp(&transaction.TransactionHandler{
		DB:     config.DB,
		Logger: logger,
		Router: apiV1,
	})

	report.StartApp(&report.ReportHandler{
		DB:     config.DB,
		Logger: logger,
		Router: apiV1,
	})

	log.Println("Server Medix BE running on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
