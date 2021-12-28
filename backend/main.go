package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"api-gofiber/pos/config"
	// "api-gofiber/pos/controllers"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	// Create a new fiber instance with custom config
	app := fiber.New(fiber.Config{
		// TIDAK DAPAT MENGHANDLE ERROR TERTENTU (TIDAK DIREKOMENDASIKAN)
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's an fiber.*Error
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			logFile, logFileError := os.OpenFile(os.Getenv("APP_LOGGER_LOCATION"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

			if logFileError != nil {
				ctx.Status(500).JSON(fiber.Map{
					"message": "Terjadi Kesalahan",
				})
			}

			logger := log.New(logFile, "Error : ", log.LstdFlags)

			logger.Println(time.Now().String())

			logger.Println(err.Error())

			fmt.Println(err.Error())

			var message string = "Terjadi Kesalahan"

			if appDebug == true {
				message = err.Error()
			}

			// Send custom error page
			err = ctx.Status(code).JSON(fiber.Map{
				"message": message,
				"code":    code,
			})

			if err != nil {
				// In case the SendFile fails
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Internal Server Error",
				})
			}

			// Return from handler
			return ctx.Status(500).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		},
	})

	// DAPAT MENGHANDLE ERROR APAPUN (RECOMENED)
	// app.Use(func(c *fiber.Ctx) error {
	// 	defer func() {
	// 		if r := recover(); r != nil {
	// 			c.Status(500).JSON(fiber.Map{
	// 				"message" : "Terjadi Kesalahan",
	// 			})
	// 		}
	// 	}()
	// 	return c.Next()
	// })

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

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(map[string]string{
			"message": "hello",
		})
	})

	api := app.Group("/api")

	v1 := api.Group("/v1") // /api/v1

	v1.Get("/test-error", func(c *fiber.Ctx) error {
		return c.SendFile("file-does-not-exist")

	})

	v1.Get("/status", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(map[string]string{
			"message": "active",
		})
	})

	// v1.Post("/register", controllers.Register)
	// v1.Post("/login", controllers.Login)

	// // JWT Middleware
	// app.Use(jwtware.New(jwtware.Config{
	// 	SigningKey: []byte("secret"),
	// }))

	// v1.Post("/refresh-token", controllers.RefreshToken)
	// v1.Post("/logout", controllers.Logout)
	// v1.Get("/me", controllers.Me)

	// v1.Put("/profil/update", controllers.UpdateProfilData)
	// v1.Post("/profil/upload", controllers.UpdateProfilPhoto)

	// v1.Get("/data", controllers.IndexData)
	// v1.POST("/data", controllers.StoreData)
	// v1.GET("/data/:id", controllers.ShowData)
	// v1.DELETE("/data/:id", controllers.DestoryData)
	// v1.PUT("/data/:id", controllers.UpdateData)

	app.Listen(os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT"))
}
