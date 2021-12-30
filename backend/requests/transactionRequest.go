package requests

import (
	"api-gofiber/pos/models"
	"errors"

	"gorm.io/gorm"
)

type DetailOrder struct {
	Qty       string `form:"qty" json:"qty" xml:"qty" validate:"required,number"`
	Price     string `form:"price" json:"price" xml:"price" validate:"required,numeric"`
	ProductID string `form:"product_id" json:"product_id" xml:"product_id" validate:"required,number"`
}

type TransactionRequest struct {
	CustomerID   string         `form:"customer_id" json:"customer_id"  xml:"customer_id" validate:"required,number"`
	DetailOrders []*DetailOrder `form:"detail_orders" json:"detail_orders" xml:"detail_orders" validate:"min=1,required,dive"`
}

func (requestData *TransactionRequest) ValidateData() error {
	return nil
}

func (requestData *TransactionRequest) ValidateExistsCustomer(database *gorm.DB) error {
	resultData := map[string]interface{}{}

	queryData := database.Model(&models.Customer{})
	queryData.Select("id")
	queryData.Where("id = ?", requestData.CustomerID)
	queryData.First(&resultData)

	if len(resultData) == 0 {
		return errors.New("Category Tidak Ditemukan")
	}

	return nil
}

func (requestData *TransactionRequest) ValidateExistsProduct(database *gorm.DB) error {
	var isValid bool = false

	for _, item := range requestData.DetailOrders {
		resultData := map[string]interface{}{}

		queryData := database.Model(&models.Customer{})
		queryData.Select("id")
		queryData.Where("id = ?", item.ProductID)
		queryData.First(&resultData)

		if len(resultData) == 0 {
			isValid = true
			break
		}
	}

	if isValid {
		return errors.New("Product Tidak Ditemukan")
	}

	return nil
}
