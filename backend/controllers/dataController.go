package controllers 

import (
	"github.com/gofiber/fiber/v2"
	"api-gofiber/test/models"
	"gorm.io/gorm"
	"strconv"
	"math"
)

func IndexData(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB);	

	page := c.Query("page", "1")
	new_page,_ := strconv.Atoi(page)

	per_page := c.Query("per_page", "10")
	new_per_page,_ := strconv.Atoi(per_page)

	search := c.Query("search","")

	var resultCount int64;

	var modelCount models.Data

	queryResultCount := database.Model(&modelCount)
		queryResultCount.Select("id","name")
		if search != "" {
			queryResultCount.Where("name LIKE ?", "%"+search+"%")		
		}
		queryResultCount.Count(&resultCount)		

    count_total_page := float64(resultCount) / float64(new_per_page)
    total_page := int(math.Ceil(count_total_page))
    limitStart := (new_page - 1) * new_per_page;
	
	result := []map[string]interface{}{}

	var queryCount models.Data
	query := database.Model(&queryCount)
		query.Select("name","id","phone")
		if search != "" {
			query.Where("name LIKE ?", "%"+search+"%")		
		}
		query.Order("id desc")
		query.Offset(limitStart)		
		query.Limit(new_per_page)		
		query.Find(&result)

	return c.Status(200).JSON(fiber.Map{
		"data" : result,
		"per_page" : new_per_page,
		"total_page" : total_page,
		"total_data" : int(resultCount),
	})
}