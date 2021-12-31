package controllers

import (
	"api-gofiber/pos/models"
	"fmt"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProductOrder struct {
	ID   int    `gorm:"column:id" json:"id"`
	Name string `gorm:"column:name" json:"name"`
}

type DetailOrder struct {
	ID        int     `gorm:"column:id" json:"id"`
	Qty       int     `gorm:"column:qty" json:"qty"`
	Price     float32 `gorm:"column:price" json:"price"`
	OrderID   *uint   `gorm:"column:order_id" json:"order_id"`
	ProductID *uint   `gorm:"column:product_id" json:"product_id"`
	// NAMA DARI STRUCT SANGAT PENGARUH SEKALI
	// ProductOrder ProductOrder `gorm:"Foreignkey:ProductID;association_foreignkey:ID" json:"product"`
	Product Product `gorm:"json:product"`
}

type Customer struct {
	ID   int    `gorm:"column:id" json:"id"`
	Name string `gorm:"column:name" json:"name"`
}

type User struct {
	ID   int    `gorm:"column:id" json:"id"`
	Name string `gorm:"column:name" json:"name"`
}

type Order struct {
	ID    int     `gorm:"column:id" json:"id"`
	Total float32 `gorm:"column:total" json:"total"`

	DetailOrder []DetailOrder `gorm:"Foreignkey:OrderID;association_foreignkey:ID" json:"detail_orders"`

	UserID *uint `gorm:"column:user_id" json:"user_id"`
	User   User  `json:"user"`

	CustomerID *uint    `gorm:"column:customer_id" json:"customer_id"`
	Customer   Customer `json:"customer"`

	DeletedAt *gorm.DeletedAt
}

func IndexOrder(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	page := c.Query("page", "1")
	new_page, _ := strconv.Atoi(page)

	per_page := c.Query("per_page", "10")
	new_per_page, _ := strconv.Atoi(per_page)

	soft_deleted := c.Query("soft_deleted", "")

	search := c.Query("search", "")

	var resultCount int64

	queryResultCount := database.Debug().Model(&models.Order{}).Preload("User").Preload("Customer").Preload("DetailOrder.Product")

	if search != "" {
		queryResultCount.Or("CAST(total as text) like ?", "%"+search+"%")
		queryResultCount.Or("EXISTS(SELECT 1 FROM users WHERE users.name = ? AND orders.user_id = users.id)", search)
	}

	if soft_deleted != "" {
		queryResultCount.Unscoped()
	}

	queryResultCount.Count(&resultCount)
	fmt.Println(resultCount)

	count_total_page := float64(resultCount) / float64(new_per_page)
	total_page := int(math.Ceil(count_total_page))
	limitStart := (new_page - 1) * new_per_page

	result := []Order{}

	query := database.Debug().Preload("User").Preload("Customer").Preload("DetailOrder.Product")

	if search != "" {
		query.Or("CAST(total as text) like ?", "%"+search+"%")
		query.Or("EXISTS(SELECT 1 FROM users WHERE users.name = ? AND orders.user_id = users.id)", search)
	}

	if soft_deleted != "" {
		query.Unscoped()
	}

	query.Order("id desc")
	query.Offset(limitStart)
	query.Limit(new_per_page)
	query.Find(&result)

	return c.Status(200).JSON(fiber.Map{
		"data":       result,
		"per_page":   new_per_page,
		"total_page": total_page,
		"total_data": resultCount,
		"page":       new_page,
	})
}

func ShowOrder(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	id, errGetParam := strconv.Atoi(c.Params("id"))

	if errGetParam != nil {
		fmt.Println(errGetParam.Error())
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	resultData := Order{}
	queryData := database.Debug().Preload("User").Preload("Customer").Preload("DetailOrder.Product")
	queryData.Where("id = ?", id)
	queryData.First(&resultData)

	if resultData.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Not Found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": resultData,
	})
}
