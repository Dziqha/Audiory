package configs

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func connectToDatabase(dsn string, retries int, delay time.Duration) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	for i := 0; i < retries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err == nil {
			return db, nil
		}
		fmt.Printf("Failed to connect to DB: %v. Retrying in %.0f seconds...\n", err, delay.Seconds())
		time.Sleep(delay)
	}
	return nil, fmt.Errorf("unable to connect to database after %d retries: %w", retries, err)
}

func Database() *gorm.DB {
	// Retrieve environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Construct DSN
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Connect to the database
	db, err := connectToDatabase(dsn, 5, 2*time.Second)
	if err != nil {
		panic(fmt.Sprintf("Database connection failed: %v", err))
	}

	fmt.Println("Connected to the PostgreSQL database successfully!")
	return db
}
