package controllers

import (
	"api-gofiber/pos/helpers"
	"api-gofiber/pos/models"
	"api-gofiber/pos/requests"
	"time"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
)

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func ValidateStruct(user requests.ValidateInterface) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func Register(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")

	users := new(requests.SignupRequest)

	if err := c.BodyParser(users); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errValidate := ValidateStruct(*users)
	if errValidate != nil {
		return c.JSON(errValidate)
	}

	hash, errorHash := helpers.HashPassword(password)

	if errorHash != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	user := models.User{
		Name:     name,
		Email:    email,
		Password: hash,
	}

	errCreateUser := database.Create(&user).Error

	if errCreateUser != nil {
		return c.Status(200).JSON(fiber.Map{
			"message": true,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": true,
	})
}

func Login(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	email := c.FormValue("email")
	password := c.FormValue("password")

	users := new(requests.SigninRequest)

	if err := c.BodyParser(users); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errValidate := ValidateStruct(*users)
	if errValidate != nil {
		return c.JSON(errValidate)
	}

	resultUser := map[string]interface{}{}

	queryUser := database.Model(&models.User{})
	queryUser.Select("id", "password", "email")
	queryUser.Where("email = ?", email)
	queryUser.First(&resultUser)

	if len(resultUser) == 0 {
		return c.Status(500).JSON(fiber.Map{
			"message": "Email tidak ditemukan",
		})
	}

	var isValidPassword bool = helpers.CheckPasswordHash(
		password,
		resultUser["password"].(string),
	)

	if isValidPassword == false {
		return c.Status(500).JSON(fiber.Map{
			"message": "Password Salah",
		})
	}

	claims := jwt.MapClaims{
		"sub": resultUser["id"],
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	access_token, errSigned := token.SignedString([]byte("secret"))

	if errSigned != nil {
		return c.Status(200).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message":      true,
		"access_token": access_token,
	})
}

func Me(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["sub"].(float64)
	return c.Status(200).JSON(fiber.Map{
		"message": true,
		"sub":     id,
	})
}

func Logout(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"message": true,
	})
}

func RefreshToken(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claimHeaders := user.Claims.(jwt.MapClaims)
	sub := claimHeaders["sub"].(float64)

	claims := jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	access_token, errSigned := token.SignedString([]byte("secret"))

	if errSigned != nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message":      true,
		"access_token": access_token,
	})
}
