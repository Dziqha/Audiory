package controller

import (
	"UserService/configs"
	"UserService/src/models"
	"UserService/src/res"
	"UserService/src/utils"
	"encoding/base64"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type AdminController struct{}

func NewAdminController() *AdminController {
	return &AdminController{}
}

func (a *AdminController) AddAdmin(c *fiber.Ctx) error {
	var req models.Admin

	if err := c.BodyParser(&req); err != nil {
		return errors.Wrap(err, "failed to parse request body")
	}

	if err := utils.ValidateStruct(&req); err != nil { 
		return errors.Wrap(err, "failed to validate request body")
	}

	// Validasi panjang password dan karakter yang diizinkan
	if len(req.Password) < 5 || len(req.Password) > 15 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Password must be between 5 and 15 characters.",
		})
	}

	// Validasi harus ada huruf
	hasLetter := regexp.MustCompile(`[A-Za-z]`).MatchString(req.Password)
	// Validasi harus ada angka
	hasDigit := regexp.MustCompile(`\d`).MatchString(req.Password)
	// Validasi harus ada karakter khusus
	hasSpecialChar := regexp.MustCompile(`[@$!%*?&]`).MatchString(req.Password)

	if !hasLetter || !hasDigit || !hasSpecialChar {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Password must contain at least one letter, one number and one special character.",
		})
	}

	encode := base64.StdEncoding.EncodeToString([]byte(req.Password))

	req.Password = string(encode)
	

	err := configs.Database().Transaction(func(tx *gorm.DB) error {
		if result := tx.Create(&req); result.Error != nil {
			return errors.Wrap(result.Error, "failed to create admin")
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "failed to process the request")
	}

	res := res.AdminResponse{
		ID: req.ID,
		Username: req.Username,
		Email: req.Email,
	}

	return c.JSON(fiber.Map{
		"message": "Admin created successfully",
		"data": res,
	})
}

func (a *AdminController) LoginAdmin(c *fiber.Ctx) error {
	var TokenSecret = os.Getenv("TOKEN_SECRET_ADMIN")
	var req models.AdminRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if err := utils.ValidateStruct(req); err != nil {
		return err
	}

	encode := base64.StdEncoding.EncodeToString([]byte(req.Password))
		
	password := string(encode)
	
		var admin models.Admin
		err := configs.Database().Transaction(func(tx *gorm.DB) error {
			result := tx.Where("username = ? AND password = ?", req.Username, password).First(&admin)
			if result.Error != nil {
				return result.Error
			}

			
			return nil
		})
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Invalid username or password",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"admin_id": admin.ID,
			"username": admin.Username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, err := token.SignedString([]byte(TokenSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	cookie := &fiber.Cookie{
		Name: "myapp_admin",
		Expires: time.Now().Add(time.Hour * 24),// Digunakan untuk mengatasi risiko CSRF 
		SameSite: fiber.CookieSameSiteStrictMode,
		HTTPOnly: true, // Cookie tidak bisa diakses dari JavaScript
		Secure: true, // Hanya dikirim melalui HTTPS
	}

	c.Cookie(cookie)
	res := res.AdminResponseLogin{
		ID:    admin.ID,
		Email: admin.Email,
		Username: admin.Username,
		Token: tokenString,
	}

	return c.JSON(fiber.Map{
		"message": "User logged in successfully",
		"user":    res,
	})
}

func (a *AdminController) GetAllAdmin(c *fiber.Ctx) error {
	var admins []models.Admin
	var finAdmin []res.AdminResponse

	err := configs.Database().Transaction(func(tx *gorm.DB) error {
		if result := tx.Find(&admins); result.Error != nil {
			return errors.Wrap(result.Error, "failed to get admin")
		}

		log.Println("Admins: ", admins)

		return nil
	})

	if len(admins) == 0 {
		return c.JSON(fiber.Map{
			"message": "No admin found",
		})
	}

	if err != nil {
		return errors.Wrap(err, "failed to process the request")
	}

	for _, resData := range admins {
		finAdmin = append(finAdmin, res.AdminResponse{
			ID: resData.ID,
			Username: resData.Username,
			Email: resData.Email,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Admin found successfully",
		"data": finAdmin,
	})
}

func (a *AdminController) GetAdminById(c *fiber.Ctx) error {
	var admin models.Admin

	id := c.Params("id")

	err := configs.Database().Transaction(func(tx *gorm.DB) error {
		if result := tx.First(&admin, id); result.Error != nil {
			return errors.Wrap(result.Error, "failed to get admin")
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "failed to process the request")
	}

	finAdmin := res.AdminResponse{
		ID: admin.ID,
		Username: admin.Username,
		Email: admin.Email,
	}

	return c.JSON(fiber.Map{
		"message": "Admin found successfully",
		"data": finAdmin,
	})


}

func (a *AdminController) UpdateAdmin(c *fiber.Ctx) error {
	var admin models.Admin
	var req models.Admin

	id := c.Params("id")

	if err := c.BodyParser(&req); err != nil {
		return errors.Wrap(err, "failed to parse request body")
	}

	err := configs.Database().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).First(&admin).Error; err != nil {
			return errors.Wrap(err, "failed to find admin by id")
		}

		if req.Username != "" {
			admin.Username = req.Username
		}
		if req.Email != "" {
			admin.Email = req.Email
		}
		if req.Password != "" {
			admin.Password = req.Password
		}

		if err := tx.Save(&admin).Error; err != nil {
			return errors.Wrap(err, "failed to update admin")
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "failed to process the request")
	}

	responseData := res.AdminResponseUpdate{
		ID: admin.ID,
		Username: admin.Username,
		Email: admin.Email,
		Password: admin.Password,
	}

	return c.JSON(fiber.Map{
		"message": "Admin updated successfully",
		"data":    responseData,
	})
}


func (a *AdminController) DeleteAdmin(c *fiber.Ctx) error {
	id := c.Params("id")
	var admin models.Admin

	err := configs.Database().Transaction(func(tx *gorm.DB) error {
		if err :=  tx.Where("id = ?", id).First(&admin).Error; err != nil {
			return errors.Wrap(err, "failed to find admin by id")
		}

		if err := tx.Delete(&admin).Error; err != nil {
			return errors.Wrap(err, "failed to delete admin")
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "failed to process the request")
	}

	return c.JSON(fiber.Map{
		"message": "Admin deleted successfully",
	})
}