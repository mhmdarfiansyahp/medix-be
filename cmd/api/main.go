package main

import (
	"log"

	"medix-be/config"
	"medix-be/internal/handler"
	"medix-be/internal/repository"
	"medix-be/internal/router"
	"medix-be/internal/service"
	"medix-be/migrations"
)

func main() {
	config.LoadEnv()

	migrations.RunMigrations()
	config.ConnectDatabase()

	db := config.DB

	obatRepo := repository.NewObatRepository(db)
	obatService := service.NewObatService(obatRepo)
	obatHandler := handler.NewObatHandler(obatService)

	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo)
	authHandler := handler.NewAuthHandler(authService)

	r := router.SetupRouter(obatHandler, authHandler)

	log.Println("Server Medix BE running on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
