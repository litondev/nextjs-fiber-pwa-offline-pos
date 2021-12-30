package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"api-gofiber/pos/config"
	"api-gofiber/pos/controllers"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtware "github.com/gofiber/jwt/v3"
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

	db, errDb := config.Database(dsn)

	if errDb != nil {
		fmt.Println(errDb)
		os.Exit(1)
	}

	debug := os.Getenv("APP_DEBUG")
	appDebug, errDebug := strconv.ParseBool(debug)
	if errDebug != nil {
		fmt.Println(errDebug)
		os.Exit(1)
	}

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		defer func() {
			if rec := recover(); rec != nil {
				logFile, logFileError := os.OpenFile(os.Getenv("APP_LOGGER_LOCATION"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

				message := rec

				if appDebug == false {
					message = "Terjadi Kesalahan"
				}

				if logFileError != nil {
					c.Status(500).JSON(fiber.Map{
						"message": message,
					})
				}

				logger := log.New(logFile, "Error : ", log.LstdFlags)
				logger.Println(time.Now().String())
				logger.Println(rec)

				fmt.Println(rec)

				c.Status(500).JSON(fiber.Map{
					"message": message,
				})
			}
		}()
		return c.Next()
	})

	app.Use(cors.New(cors.Config{
		Next:             nil,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	}))

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("DB", db)
		c.Locals("DEBUG", debug)
		return c.Next()
	})

	app.Static("/assets", "./assets")

	api := app.Group("/api")

	v1 := api.Group("/v1")

	v1.Get("/test-error", func(c *fiber.Ctx) error {
		a := 1
		b := 0
		zero := a / b
		return c.JSON(fiber.Map{
			"zero": zero,
		})

		// return c.SendFile("file-does-not-exist")
	})

	v1.Get("/status", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(map[string]string{
			"message": "active",
		})
	})

	v1.Post("/register", controllers.Register)
	v1.Post("/login", controllers.Login)

	// // JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("secret"),
	}))

	v1.Post("/refresh-token", controllers.RefreshToken)
	v1.Get("/me", controllers.Me)

	v1.Get("/category", controllers.IndexCategory)
	v1.Post("/category", controllers.StoreCategory)
	v1.Get("/category/:id", controllers.ShowCategory)
	v1.Delete("/category/:id", controllers.DestroyCategory)
	v1.Put("/category/:id", controllers.UpdateCategory)

	v1.Get("/customer", controllers.IndexCustomer)
	v1.Post("/customer", controllers.StoreCustomer)
	v1.Get("/customer/:id", controllers.ShowCustomer)
	v1.Delete("/customer/:id", controllers.DestroyCustomer)
	v1.Put("/customer/:id", controllers.UpdateCustomer)

	v1.Get("/user", controllers.IndexUser)
	v1.Post("/user", controllers.StoreUser)
	v1.Get("/user/:id", controllers.ShowUser)
	v1.Delete("/user/:id", controllers.DestroyUser)
	v1.Put("/user/:id", controllers.UpdateUser)

	v1.Get("/product", controllers.IndexProduct)
	v1.Get("/product/code", controllers.GetCodeProduct)
	v1.Post("/product", controllers.StoreProduct)
	v1.Get("/product/:id", controllers.ShowProduct)
	v1.Delete("/product/:id", controllers.DestroyProduct)
	v1.Put("/product/:id", controllers.UpdateProduct)

	app.Listen(os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT"))
}
