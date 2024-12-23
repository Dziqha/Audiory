package controller

import (
	"UserService/configs"
	"UserService/src/models"
	"UserService/src/res"
	"UserService/src/utils"
	"crypto/rand"
	"encoding/hex"
	_ "fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SubscriptionController struct{}

func NewSubscriptionController() *SubscriptionController {
	return &SubscriptionController{}
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

func (s *SubscriptionController) AddSubscription(c *fiber.Ctx) error {
	var reqData models.SubscriptionRequest

	if err := c.BodyParser(&reqData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse request body",
		})
	}

	if err := utils.ValidateStruct(&reqData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var user models.User
	if err := configs.Database().First(&user, reqData.UserID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	var existingSubscription models.Subscription
	if err := configs.Database().Where("user_id = ? AND is_active = ?", user.ID, true).First(&existingSubscription).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User already has an active subscription",
		})
	}


	req := models.Subscription{
		SubscriptionStart: time.Now(),
		SubscriptionType:  reqData.SubscriptionType, 
		UserID:            user.ID,
		IsActive:          true,
	}
	var price int

	switch req.SubscriptionType {
	case "day":
		req.SubscriptionEnd = req.SubscriptionStart.AddDate(0, 0, 1) 
		price = 10000 
	case "month":
		req.SubscriptionEnd = req.SubscriptionStart.AddDate(0, 1, 0) 
		price = 250000 
	case "year":
		req.SubscriptionEnd = req.SubscriptionStart.AddDate(1, 0, 0) 
		price = 2500000 
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid subscription type",
		})
	}
	var wg sync.WaitGroup
	paymentChan := make(chan *res.PaymentResponse, 1) // menyimpan satu respons pembayaran
	errorChan := make(chan error, 1) // menyimpan satu kesalahan

	wg.Add(1)
	go func() {
		defer wg.Done()
		paymentResponse, err := MidtransPayment(req.UserID, generateRandomString(10), price, req.SubscriptionType)
		if err != nil {
			errorChan <- err
			return
		}
		paymentChan <- paymentResponse // Mengirimkan data ke paymentChan
	}()

	wg.Wait()
	close(paymentChan)
	close(errorChan)

	if err := <-errorChan; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create payment transaction",
			"error":   err.Error(),
		})
	}

	paymentResponse := <-paymentChan  // Menerima data dari paymentChan

	req.SubscriptionToken = paymentResponse.Token



	err := configs.Database().Transaction(func(tx *gorm.DB) error {
		if result := tx.Create(&req); result.Error != nil {
			return result.Error
		}
		return nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	response := res.SubscriptionResponse{
		SubscriptionStart: req.SubscriptionStart,
		SubscriptionEnd:   req.SubscriptionEnd,
		SubscriptionToken: req.SubscriptionToken,
		SubscriptionType:  req.SubscriptionType,
		IsActive:          req.IsActive,
		UserID:            req.UserID,
		PaymentToken:      paymentResponse.Token,
		PaymentRedirectURL: paymentResponse.RedirectURL,
	}

	return c.JSON(fiber.Map{
		"message": "Subscription added successfully",
		"data":    response,
	})
}

func (s *SubscriptionController) RemoveExpiredSubscriptions() error {
	var subscriptions []models.Subscription
	currentTime := time.Now()

	// Cari langganan yang sudah kadaluarsa dan masih aktif
	if err := configs.Database().Where("subscription_end < ? AND is_active = ?", currentTime, true).Find(&subscriptions).Error; err != nil {
		return err
	}

	// Set is_active menjadi false untuk langganan yang kadaluarsa
	for _, subscription := range subscriptions {
		subscription.IsActive = false
		if err := configs.Database().Save(&subscription).Error; err != nil {
			return err
		}
	}

	return nil
}
