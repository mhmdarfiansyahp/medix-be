package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadEnv() {
	viper.SetConfigName("config")   
	viper.SetConfigType("yaml")     
	viper.AddConfigPath("./config") 
	viper.AddConfigPath(".")      

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error loading config.yaml file: %v", err)
	}

	// Ensure JWT_SECRET is available for packages using os.Getenv
	if os.Getenv("JWT_SECRET") == "" {
		if secret := viper.GetString("jwt.secret"); secret != "" {
			os.Setenv("JWT_SECRET", secret)
		}
	}
}

func ConnectDatabase() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		viper.GetString("database.host"),
		viper.GetString("database.username"),
		viper.GetString("database.password"),
		viper.GetString("database.name"),
		viper.GetString("database.port"),
		viper.GetString("database.ssl"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Create uuid-ossp extension if you plan to use UUID in PostgreSQL
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		log.Printf("Warning creating uuid-ossp extension: %v", err)
	}

	DB = db
	fmt.Println("Successfully connected to Medix database!")
}
