package main

import (
	"api-gofiber/pos/config"
	"api-gofiber/pos/helpers"
	"api-gofiber/pos/models"
	"fmt"
	"os"
	"strconv"
	"time"

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
		// USER
		for i := 0; i < 15; i++ {
			hash, _ := helpers.HashPassword("password")
			var name string = "user" + strconv.Itoa(i)

			user := models.User{
				Name:      name,
				Email:     name + "@gmail.com",
				Password:  hash,
				DeletedAt: nil,
			}

			result := database.Create(&user)

			if result.Error == nil {
				fmt.Println("Success Memasukan " + strconv.Itoa(int(result.RowsAffected)) + " Data")
			} else {
				fmt.Print(result.Error)
				os.Exit(1)
			}
		}

		// CATEGORY
		for i := 0; i < 15; i++ {
			var name string = "category" + strconv.Itoa(i)
			var value string = "Des"
			var description *string = &value

			if i == 2 {
				description = nil
			} else if i == 5 {
				description = nil
			}

			data := models.Category{
				Name:        name,
				Description: description,
				DeletedAt:   nil,
			}

			result := database.Create(&data)

			if result.Error == nil {
				fmt.Println("Success Memasukan " + strconv.Itoa(int(result.RowsAffected)) + " Data")
			} else {
				fmt.Print(result.Error)
				os.Exit(1)
			}
		}

		// PRODUCT
		for i := 0; i < 15; i++ {
			var name string = "product" + strconv.Itoa(i)

			now := time.Now()
			var code string = strconv.Itoa(now.Year()) + strconv.Itoa(int(now.Month())) + strconv.Itoa(now.Day()) + "-" + strconv.Itoa(now.Nanosecond())

			var valueDes string = "Des"
			var description *string = &valueDes
			var stock int = 100
			var price float32 = 100000.00
			var categoryId uint = 1

			if i == 2 || i == 5 || i == 8 {
				description = nil
				stock = 200
				price = 200000.00
				categoryId = 2
			}

			data := models.Product{
				Name:        name,
				Code:        code,
				Description: description,
				Stock:       stock,
				Price:       price,
				CategoryID:  &categoryId,
				Photo:       nil,
				DeletedAt:   nil,
			}

			result := database.Create(&data)

			if result.Error == nil {
				fmt.Println("Success Memasukan " + strconv.Itoa(int(result.RowsAffected)) + " Data")
			} else {
				fmt.Print(result.Error)
				os.Exit(1)
			}
		}

		// CUSTOMER
		for i := 0; i < 15; i++ {
			var name string = "customer" + strconv.Itoa(i)
			var gmail string = name + "@gmail.com"
			var email *string = &gmail

			customer := models.Customer{
				Name:      name,
				Email:     email,
				Address:   nil,
				Phone:     nil,
				DeletedAt: nil,
			}

			result := database.Create(&customer)

			if result.Error == nil {
				fmt.Println("Success Memasukan " + strconv.Itoa(int(result.RowsAffected)) + " Data")
			} else {
				fmt.Print(result.Error)
				os.Exit(1)
			}
		}

		// ORDER
		for i := 0; i < 15; i++ {
			var Id uint = 1
			pId := &Id
			order := models.Order{
				Total:      100000,
				UserID:     pId,
				CustomerID: pId,
			}

			result := database.Create(&order)

			if result.Error == nil {
				fmt.Println("Success Memasukan " + strconv.Itoa(int(result.RowsAffected)) + " Data")
			} else {
				fmt.Print(result.Error)
				os.Exit(1)
			}
		}

		// DETAIL ORDER
		for i := 0; i < 15; i++ {
			var Id uint = 1
			pId := &Id

			detailOrder := models.DetailOrder{
				Qty:       10,
				Price:     10000,
				OrderID:   pId,
				ProductID: pId,
			}

			result := database.Create(&detailOrder)

			if result.Error == nil {
				fmt.Println("Success Memasukan " + strconv.Itoa(int(result.RowsAffected)) + " Data")
			} else {
				fmt.Print(result.Error)
				os.Exit(1)
			}
		}

		fmt.Println("Success Seed")
	}
}
