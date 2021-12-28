package main

import (
	"api-gofiber/pos/config"
	"api-gofiber/pos/models"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	errEnv := godotenv.Load()

	if errEnv != nil {
		fmt.Println(errEnv)
		os.Exit(1)
	}

	var dsn map[string]string = map[string]string{
		"DB_HOST":     os.Getenv("DB_HOST"),
		"DB_PORT":     os.Getenv("DB_PORT"),
		"DB_NAME":     os.Getenv("DB_NAME"),
		"DB_USER":     os.Getenv("DB_USER"),
		"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
	}

	database, errDb := config.Database(dsn)

	if errDb == nil {
		database.Migrator().DropTable(&models.User{})
		database.Migrator().DropTable(&models.Category{})
		database.Migrator().DropTable(&models.Product{})
		database.Migrator().DropTable(&models.Customer{})
		database.Migrator().DropTable(&models.Order{})
		database.Migrator().DropTable(&models.DetailOrder{})

		database.AutoMigrate(&models.User{}, &models.Category{}, &models.Product{}, &models.Customer{}, &models.Order{}, &models.DetailOrder{})

		fmt.Println("Success Migrate")
	}
}
