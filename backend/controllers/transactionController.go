package controllers

import (
	"api-gofiber/pos/models"
	"api-gofiber/pos/requests"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

func Transaction(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	transaction := new(requests.TransactionRequest)

	if errParser := c.BodyParser(transaction); errParser != nil {
		fmt.Println(errParser.Error())

		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	errValidate := requests.ValidateStruct(transaction)
	if errValidate != nil {
		return c.Status(422).JSON(fiber.Map{
			"errors": errValidate,
		})
	}

	errValidateExistsCustomer := transaction.ValidateExistsCustomer(database)

	if errValidateExistsCustomer != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": errValidateExistsCustomer.Error(),
		})
	}

	errValidateExistsProduct := transaction.ValidateExistsProduct(database)

	if errValidateExistsProduct != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": errValidateExistsProduct.Error(),
		})
	}

	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	id := claims["sub"].(float64)

	var realId uint = uint(id)
	newRealId := &realId

	pointerCustomerId := transaction.CustomerID

	tx := database.Begin()

	order := models.Order{
		UserID:     newRealId,
		CustomerID: pointerCustomerId,
	}

	errCreateOrder := database.Create(&order).Error

	if errCreateOrder != nil {
		fmt.Println(errCreateOrder.Error())
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	var total float32 = 0

	for _, item := range transaction.DetailOrders {
		// qty, _ := strconv.ParseInt(item.Qty, 10, 32)
		realQty := item.Qty
		// price, _ := strconv.ParseFloat(item.Price, 10)
		realPrice := item.Price
		// productId, _ := strconv.ParseInt(item.ProductID, 10, 32)
		// realProductId := item.ProductID
		pointerProductId := item.ProductID
		pointerOrderId := &order.ID
		detailOrder := models.DetailOrder{
			Qty:       realQty,
			Price:     realPrice,
			ProductID: pointerProductId,
			OrderID:   pointerOrderId,
		}

		errCreateDetailOrder := database.Create(&detailOrder).Error

		if errCreateDetailOrder != nil {
			fmt.Println(errCreateDetailOrder.Error())
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{
				"message": "Terjadi Kesalahan",
			})
		}

		total += realPrice
	}

	order.Total = total

	queryData := database.Model(&models.Order{})
	queryData.Select("total")
	queryData.Where("id = ?", order.ID)

	errUpdateOrder := queryData.Updates(&order).Error

	if errUpdateOrder != nil {
		fmt.Println(errUpdateOrder.Error())
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	tx.Commit()
	return c.Status(200).JSON(fiber.Map{
		"message": true,
	})
}
