package controllers 

import (	
	"os"
	"fmt"	
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"api-gofiber/test/helpers"
	"api-gofiber/test/models"
	// "api-gofiber/test/requests"
	"path/filepath"
	"github.com/disintegration/imaging"
	"gorm.io/gorm"
)

func UpdateProfilData(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB);	
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	name := c.FormValue("name");
	email := c.FormValue("email");
	password := c.FormValue("password");
	password_confirm := c.FormValue("password_confirm")

	var ID uint = uint(claims["sub"].(float64))

	resultSearchEmail := map[string]interface{}{}
	
	querySearchEmail := database.Model(&models.User{})
		querySearchEmail.Where("email = ?",email)
		querySearchEmail.Not("id = ?",ID)
		querySearchEmail.First(&resultSearchEmail)

	if len(resultSearchEmail) > 0 {
		return c.Status(500).JSON(fiber.Map{
			"message": "Email telah terpakai",
		})
	}

	resultUser := map[string]interface{}{}

	queryUser := database.Model(&models.User{})
		queryUser.Where("id = ?",ID)
		queryUser.First(&resultUser)

	var isValidPassword bool = helpers.CheckPasswordHash(
		password_confirm,
		resultUser["password"].(string),
	)

	if isValidPassword == false {
		return c.Status(500).JSON(fiber.Map{
			"message": "Password Konfirmasi Tidak Valid",
		})		
	}

	updateUser := &models.User{
		Email : email,
		Name : name,
	}

	if(password != ""){
		hash, errHash := helpers.HashPassword(password)

		if(errHash != nil){
			fmt.Println(errHash.Error())		
			return c.Status(500).JSON(fiber.Map{
				"message" : "Terjadi Kesalahan",
			})			
		}

		updateUser.Password = hash
	}
		
	queryUpdateUser := database.Model(&models.User{})
		queryUpdateUser.Where("id = ?",ID)
		queryUpdateUser.Updates(&updateUser);

	if queryUpdateUser.Error != nil {
		fmt.Println(queryUpdateUser.Error)
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message" : true,
	})
}

func UpdateProfilPhoto(c *fiber.Ctx) error {
	database := c.Locals("DB").(*gorm.DB);	

	file, errGetFile := c.FormFile("photo")

	if errGetFile != nil {
		fmt.Println(errGetFile.Error())		
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})		
	}

	if(false == (file.Header["Content-Type"][0] != "image/jpeg" || file.Header["Content-Type"][0] != "image/png" || file.Header["Content-Type"][0] != "image/jpg")){
		fmt.Println(file.Header["Content-Type"][0])
		return c.Status(500).JSON(fiber.Map{
			"message": "Gambar tidak valid",
		})	 	
	}	

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	var ID uint = uint(claims["sub"].(float64))

    var stringID string = strconv.FormatUint(uint64(ID),10)

	filename := stringID + "-" + file.Filename;
	/* FOLDER USER HARUS ADA */
	pathname := filepath.Base("") + "/assets/users/" + filename

	if _,errFileExists := os.Stat(pathname); errFileExists == nil {
		errRemoveFile := os.Remove(pathname)
		if errRemoveFile != nil {
			fmt.Println(errRemoveFile.Error())
			return c.Status(500).JSON(fiber.Map{
				"message" : "Terjadi Kesalahan",
			})
		}
	}

	if errUploadFile := c.SaveFile(file, pathname); errUploadFile != nil {				
		fmt.Println(errUploadFile.Error())
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	openFile, errOpenFile := imaging.Open(pathname)

	if(errOpenFile != nil){
		fmt.Println(errOpenFile.Error())
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}

	resizeFile := imaging.Resize(openFile, 128, 128, imaging.Lanczos)
	
	errRessizeFile := imaging.Save(resizeFile, pathname)

	if(errRessizeFile != nil){
		fmt.Println(errRessizeFile.Error())
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})	
	}

	queryUpdateUser := database.Model(&models.User{})
		queryUpdateUser.Select("photo")
		queryUpdateUser.Where("id = ?",ID)
		queryUpdateUser.Update("photo",filename)

	if queryUpdateUser.Error != nil {
		fmt.Println(queryUpdateUser.Error)
		return c.Status(500).JSON(fiber.Map{
			"message": "Terjadi Kesalahan",
		})
	}
		
	return c.Status(200).JSON(fiber.Map{
		"message": true,
	})			
}

