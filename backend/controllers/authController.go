package controllers

import (
	"api-gofiber/pos/helpers"
	"api-gofiber/pos/models"
	"api-gofiber/pos/requests"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"

	"github.com/golang-jwt/jwt/v4"
)

func Register(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	users := new(requests.SignupRequest)

	if errParser := c.BodyParser(users); errParser != nil {
		fmt.Println(errParser.Error())

		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	users.ValidateData()

	errValidate := requests.ValidateStruct(users)
	if errValidate != nil {
		return c.Status(422).JSON(fiber.Map{
			"errors": errValidate,
		})
	}

	hash, errHash := helpers.HashPassword(users.Password)

	if errHash != nil {
		fmt.Println(errHash.Error())

		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	user := models.User{
		Name:     users.Name,
		Email:    users.Email,
		Password: hash,
	}

	errCreateUser := database.Create(&user).Error

	if errCreateUser != nil {
		fmt.Println(errCreateUser.Error())

		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": true,
	})
}

func Login(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	users := new(requests.SigninRequest)

	if errParser := c.BodyParser(users); errParser != nil {
		fmt.Println(errParser.Error())

		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	users.ValidateData()

	errValidate := requests.ValidateStruct(users)
	if errValidate != nil {
		return c.Status(422).JSON(fiber.Map{
			"errors": errValidate,
		})
	}

	resultUser := map[string]interface{}{}

	queryUser := database.Model(&models.User{})
	queryUser.Select("id", "password", "email")
	queryUser.Where("email = ?", users.Email)
	queryUser.First(&resultUser)

	if len(resultUser) == 0 {
		return c.Status(500).JSON(fiber.Map{
			"message": "Email tidak ditemukan",
		})
	}

	isValidPassword := helpers.CheckPasswordHash(
		users.Password,
		resultUser["password"].(string),
	)

	if isValidPassword == false {
		return c.Status(500).JSON(fiber.Map{
			"message": "Password Salah",
		})
	}

	claims := jwt.MapClaims{
		"sub": resultUser["id"],
		"exp": time.Now().Add(time.Minute * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	access_token, errSigned := token.SignedString([]byte("secret"))

	if errSigned != nil {
		fmt.Println(errSigned.Error())

		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"access_token": access_token,
	})
}

func Me(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	id := claims["sub"].(float64)

	return c.Status(200).JSON(fiber.Map{
		"sub":        id,
		"user-token": userToken,
	})
}

func RefreshToken(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)
	claimHeaders := userToken.Claims.(jwt.MapClaims)
	sub := claimHeaders["sub"].(float64)

	claims := jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(time.Minute * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	access_token, errSigned := token.SignedString([]byte("secret"))

	if errSigned != nil {
		fmt.Println(errSigned.Error())

		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"access_token": access_token,
	})
}
