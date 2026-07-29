package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadEnv() {
	viper.SetConfigName("config")   // Nama file tanpa ekstensi (.yaml)
	viper.SetConfigType("yaml")     // Tipe file konfigurasi
	viper.AddConfigPath("./config") // Mencari di dalam folder config/
	viper.AddConfigPath(".")        // Alternatif jika config.yaml ditaruh di root folder

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error loading config.yaml file: %v", err)
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

	// Buat ekstensi UUID jika kamu berencana menggunakan UUID di PostgreSQL
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		log.Printf("Warning creating uuid-ossp extension: %v", err)
	}

	DB = db
	fmt.Println("Berhasil terhubung ke database Medix!")
}
