package database

import (
	"fmt"
	"log"
	"os"

	"github.com/mnadev/limestone/internal/application/domain/entity"
	"gorm.io/gorm/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDatabase() *gorm.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host,
		port,
		dbUser,
		dbName,
		password,
	)
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	log.Println("Database connected successfully.")
	err = DB.AutoMigrate(entity.Adhan{})
	if err != nil {
		return nil
	}
	err = DB.AutoMigrate(entity.Event{})
	if err != nil {
		return nil
	}
	err = DB.AutoMigrate(entity.Masjid{})
	if err != nil {
		return nil
	}
	err = DB.AutoMigrate(entity.User{})
	if err != nil {
		return nil
	}
	err = DB.AutoMigrate(entity.NikkahProfile{})
	if err != nil {
		log.Printf("Failed to migrate NikkahProfile: %v", err)
		return nil
	}
	err = DB.AutoMigrate(entity.NikkahLike{})
	if err != nil {
		log.Printf("Failed to migrate NikkahLike: %v", err)
		return nil
	}
	err = DB.AutoMigrate(entity.NikkahMatch{})
	if err != nil {
		return nil
	}
	err = DB.AutoMigrate(entity.RevertProfile{})
	if err != nil {
		return nil
	}
	err = DB.AutoMigrate(entity.RevertMatch{})
	if err != nil {
		return nil
	}
	return DB
}

func SetupDatabaseTesting() *gorm.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("TEST_DB_NAME")
	dbUser := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host,
		port,
		dbUser,
		dbName,
		password,
	)
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	log.Println("Database Testing connected successfully.")
	err = DB.AutoMigrate(entity.Adhan{})
	if err != nil {
		return nil
	}
	err = DB.AutoMigrate(entity.Event{})
	if err != nil {
		return nil
	}
	err = DB.AutoMigrate(entity.Masjid{})
	if err != nil {
		return nil
	}
	err = DB.AutoMigrate(entity.User{})
	if err != nil {
		return nil
	}
	err = DB.AutoMigrate(entity.NikkahProfile{})
	if err != nil {
		log.Printf("Failed to migrate NikkahProfile: %v", err)
		return nil
	}
	err = DB.AutoMigrate(entity.NikkahLike{})
	if err != nil {
		log.Printf("Failed to migrate NikkahLike: %v", err)
		return nil
	}
	err = DB.AutoMigrate(entity.NikkahMatch{})
	if err != nil {
		log.Printf("Failed to migrate NikkahMatch: %v", err)
		return nil
	}
	err = DB.AutoMigrate(entity.RevertProfile{})
	if err != nil {
		return nil
	}
	err = DB.AutoMigrate(entity.RevertMatch{})
	if err != nil {
		return nil
	}
	return DB
}
