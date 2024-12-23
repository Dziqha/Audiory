package controller

import (
	"UserService/configs"
	"UserService/src/models"
	"UserService/src/res"
	"UserService/src/utils"
	"encoding/base64"
	"os"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (u *UserController) Register(c *fiber.Ctx) error {
	var req models.User

	if err := c.BodyParser(&req); nil != err {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := utils.ValidateStruct(&req); nil != err {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
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
	

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: encode,
	}

	err := configs.Database().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; nil != err {
			return err
		}

		return nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	res := res.UserRes{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	return c.JSON(fiber.Map{
		"message": "User created successfully",
		"user":    res,
	})
}

func (u *UserController) Login(c *fiber.Ctx) error {
	var TokenSecret = os.Getenv("TOKEN_SECRET_USER")
	var req models.UserRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if err := utils.ValidateStruct(req); err != nil {
		return err
	}

	encode := base64.StdEncoding.EncodeToString([]byte(req.Password))
		
	password := string(encode)
	
		var user models.User
		err := configs.Database().Transaction(func(tx *gorm.DB) error {
			result := tx.Where("username = ? AND password = ?", req.Username, password).First(&user)
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
			"user_id": user.ID,
			"username": user.Username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, err := token.SignedString([]byte(TokenSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	cookie := &fiber.Cookie{
		Name: "myapp_user",
		Expires: time.Now().Add(time.Hour * 24),// Digunakan untuk mengatasi risiko CSRF 
		SameSite: fiber.CookieSameSiteStrictMode,
		HTTPOnly: true, // Cookie tidak bisa diakses dari JavaScript
		Secure: true, // Hanya dikirim melalui HTTPS
	}

	c.Cookie(cookie)
	res := res.UserLoginRes{
		ID:    user.ID,
		Email: user.Email,
		Username: user.Username,
		Token: tokenString,
	}

	return c.JSON(fiber.Map{
		"message": "User logged in successfully",
		"user":    res,
	})
}