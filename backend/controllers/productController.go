package controllers

import (
	"errors"
	"fmt"
	"math"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"api-gofiber/pos/models"
	"api-gofiber/pos/requests"

	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Category struct {
	ID   int    `gorm:"column:id" json:"id"`
	Name string `gorm:"column:name" json:"name"`
}

type Product struct {
	ID          int      `gorm:"column:id" json:"id"`
	Name        string   `gorm:"column:name" json:"name"`
	Code        string   `gorm:"column:code" json:"code"`
	Description *string  `gorm:"column:description" json:"description"`
	Stock       int      `gorm:"column:stock" json:"stock"`
	Price       float32  `gorm:"column:price" json:"price"`
	Photo       *string  `gorm:"column:photo" json:"photo"`
	CategoryID  *uint    `gorm:"column:category_id" json:"category_id"`
	Category    Category `json:"category"`
	DeletedAt   *gorm.DeletedAt
}

func IndexProduct(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	page := c.Query("page", "1")
	new_page, _ := strconv.Atoi(page)

	per_page := c.Query("per_page", "10")
	new_per_page, _ := strconv.Atoi(per_page)

	soft_deleted := c.Query("soft_deleted", "")

	search := c.Query("search", "")

	var resultCount int64

	var modelCount models.Product

	queryResultCount := database.Model(&modelCount)
	queryResultCount.Select("id")

	if search != "" {
		queryResultCount.Where(
			"(name LIKE ? OR CAST(price as text) LIKE ?) OR code LIKE ?",
			"%"+search+"%",
			"%"+search+"%",
			"%"+search+"%",
		)
	}

	if soft_deleted != "" {
		queryResultCount.Unscoped()
	}

	queryResultCount.Count(&resultCount)

	count_total_page := float64(resultCount) / float64(new_per_page)
	total_page := int(math.Ceil(count_total_page))
	limitStart := (new_page - 1) * new_per_page

	result := []Product{}

	query := database.Debug().Preload("Category")
	if search != "" {
		query.Where(
			"(name LIKE ? OR CAST(price as text) LIKE ?) OR code LIKE ?",
			"%"+search+"%",
			"%"+search+"%",
			"%"+search+"%",
		)
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

func GetAllProduct(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	result := []map[string]interface{}{}

	database.Model(&models.Product{}).Find(&result)

	return c.Status(200).JSON(fiber.Map{
		"data": result,
	})
}

func GetCodeProduct(c *fiber.Ctx) error {
	now := time.Now()
	code := strconv.Itoa(now.Year()) + strconv.Itoa(int(now.Month())) + strconv.Itoa(now.Day()) + "-" + strconv.Itoa(now.Nanosecond())

	return c.Status(200).JSON(fiber.Map{
		"code": code,
	})
}

func UploadProduct(file *multipart.FileHeader, c *fiber.Ctx) (*string, error) {
	time := time.Now()

	filename := strconv.Itoa(time.Nanosecond()) + "-" + file.Filename

	pathname := filepath.Base("./") + "/assets/images/products/" + filename

	// if _, errFileExists := os.Stat(pathname); errFileExists == nil {
	// 	errRemoveFile := os.Remove(pathname)
	// 	if errRemoveFile != nil {
	// 		fmt.Println(errRemoveFile.Error())
	// 		return nil, errors.New("Terjadi Kesalahan Pada Saat Upload")
	// 	}
	// }

	if errUploadFile := c.SaveFile(file, pathname); errUploadFile != nil {
		fmt.Println(errUploadFile.Error())
		return nil, errors.New("Terjadi Kesalahan Pada Saat Upload")
	}

	openFile, errOpenFile := imaging.Open(pathname)

	if errOpenFile != nil {
		fmt.Println(errOpenFile.Error())
		return nil, errors.New("Terjadi Kesalahan Pada Saat Upload")
	}

	resizeFile := imaging.Resize(openFile, 128, 128, imaging.Lanczos)

	errRessizeFile := imaging.Save(resizeFile, pathname)

	if errRessizeFile != nil {
		fmt.Println(errRessizeFile.Error())
		return nil, errors.New("Terjadi Kesalahan Pada Saat Upload")
	}

	return &filename, nil
}

func StoreProduct(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	product := new(requests.ProductRequest)

	if errParser := c.BodyParser(product); errParser != nil {
		fmt.Println(errParser.Error())

		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	product.ValidateData()

	errorCategoryValidation := product.ValidateExistsCategory(database)
	if errorCategoryValidation != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": errorCategoryValidation.Error(),
		})
	}

	errValidate := requests.ValidateStruct(product)
	if errValidate != nil {
		return c.Status(422).JSON(fiber.Map{
			"errors": errValidate,
		})
	}

	file, _ := c.FormFile("photo")
	var filename *string = nil
	var errorUpload error = nil

	if file != nil {
		errorPhotoValidation := product.ValidatePhoto(file)

		if errorPhotoValidation != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": errorPhotoValidation.Error(),
			})
		}

		filename, errorUpload = UploadProduct(file, c)

		if errorUpload != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": errorUpload.Error(),
			})
		}
	}

	price := product.Price
	stock := product.Stock
	realCategoryData := product.CategoryID

	productModel := models.Product{
		Code:        product.Code,
		Name:        product.Name,
		Description: product.Description,
		Stock:       int(stock),
		Price:       float32(price),
		CategoryID:  realCategoryData,
		Photo:       filename,
	}

	errCreateProduct := database.Create(&productModel).Error

	if errCreateProduct != nil {
		fmt.Println(errCreateProduct.Error())

		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": true,
	})
}

func ShowProduct(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	id, errGetParam := strconv.Atoi(c.Params("id"))

	if errGetParam != nil {
		fmt.Println(errGetParam.Error())
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	resultData := Product{}
	queryData := database.Preload("Category")
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

func UpdateProduct(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	id, errGetParam := strconv.Atoi(c.Params("id"))

	if errGetParam != nil {
		fmt.Println(errGetParam.Error())
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	product := new(requests.ProductRequest)

	if errParser := c.BodyParser(product); errParser != nil {
		fmt.Println(errParser.Error())
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	product.ValidateData()

	errorCategoryValidation := product.ValidateExistsCategory(database)
	if errorCategoryValidation != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": errorCategoryValidation.Error(),
		})
	}

	errValidate := requests.ValidateStruct(product)
	if errValidate != nil {
		return c.Status(422).JSON(fiber.Map{
			"errors": errValidate,
		})
	}

	file, _ := c.FormFile("photo")
	var filename *string = nil
	var errorUpload error = nil

	if file != nil {
		errorPhotoValidation := product.ValidatePhoto(file)

		if errorPhotoValidation != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": errorPhotoValidation.Error(),
			})
		}

		filename, errorUpload = UploadProduct(file, c)

		if errorUpload != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": errorUpload.Error(),
			})
		}
	}

	queryData := database.Model(&models.Product{})
	queryData.Select("name", "code", "description", "stock", "price", "photo", "category_id")
	queryData.Where("id = ?", id)

	resultData := map[string]interface{}{}

	queryData.First(&resultData)

	if len(resultData) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Not Found",
		})
	}

	photo := resultData["photo"]
	if file != nil && photo != nil {
		pathname := filepath.Base("./") + "/assets/images/products/" + *photo.(*string)

		if _, errFileExists := os.Stat(pathname); errFileExists == nil {
			errRemoveFile := os.Remove(pathname)
			if errRemoveFile != nil {
				fmt.Println(errRemoveFile.Error())
				return c.Status(500).JSON(fiber.Map{
					"message": "Terjadi Kesalahan",
				})
			}
		}
	}

	price := product.Price
	stock := product.Stock
	realCategoryData := product.CategoryID

	queryData.Updates(&models.Product{
		Code:        product.Code,
		Name:        product.Name,
		Description: product.Description,
		Stock:       int(stock),
		Price:       float32(price),
		CategoryID:  realCategoryData,
		Photo:       filename,
	})

	if queryData.Error != nil {
		fmt.Println(queryData.Error)

		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": true,
	})
}

func DestroyProduct(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB)

	id, errGetParam := strconv.Atoi(c.Params("id"))

	if errGetParam != nil {
		fmt.Println(errGetParam.Error())
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	queryData := database.Model(&models.Product{})
	queryData.Where("id = ?", id)

	resultData := map[string]interface{}{}

	queryData.First(&resultData)

	if len(resultData) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Not Found",
		})
	}

	queryData.Delete(&models.Product{})

	if queryData.Error != nil {
		fmt.Println(queryData.Error)
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": true,
	})
}
