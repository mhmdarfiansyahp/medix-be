// File: migrations/migration.go
package migrations

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"
)

func RunMigrations() {
	dbUser := viper.GetString("database.username")
	dbPass := viper.GetString("database.password")
	dbHost := viper.GetString("database.host")
	dbPort := viper.GetString("database.port")
	dbName := viper.GetString("database.name")
	sslMode := viper.GetString("database.ssl")

	migrationPath := viper.GetString("migration_path")
	if migrationPath == "" {
		migrationPath = "file://migrations/data"
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPass, dbHost, dbPort, dbName, sslMode)

	m, err := migrate.New(
		migrationPath,
		dsn,
	)
	if err != nil {
		log.Fatalf("Gagal inisialisasi instance migrasi: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Gagal menjalankan migrasi UP: %v", err)
	}

	fmt.Println("Migrasi database berhasil dijalankan!")
}
