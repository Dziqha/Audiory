package controller

import (
	"UserService/configs"
	"UserService/src/models"
	"UserService/src/res"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type RecomendationController struct{}

func NewRecomendationController() *RecomendationController {
	return &RecomendationController{}
}

// AddRecommendation handles adding a new recommendation based on favorite songs
func (r *RecomendationController) AddRecommendation(c *fiber.Ctx) error {
	userIDStr := c.Params("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var favoriteSongs []models.FavoriteSong
	if err := configs.Database().Where("user_id = ?", userID).Find(&favoriteSongs).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch favorite songs"})
	}

	if len(favoriteSongs) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "No favorite songs found for this user"})
	}

	genreCount := make(map[int]int)
	for _, favorite := range favoriteSongs {
		genreCount[favorite.GenreID]++
	}

	var mostLikedGenreID, maxCount int
	for genreID, count := range genreCount {
		if count > maxCount {
			maxCount = count
			mostLikedGenreID = genreID
			log.Println("mostLikedGenreID:", mostLikedGenreID)
		}
	}

	if len(favoriteSongs) == 1 {
		log.Println("Hanya ada satu lagu favorit, menggunakan lagu yang sama untuk rekomendasi")
		recommendation := models.Recomendation{
			SongID:     favoriteSongs[0].SongID,  
			UserID:     userID,                 
			GenreID:    mostLikedGenreID,
			Created_at: time.Now(),
			Reason:     fmt.Sprintf("Recommended based on your favorite song with ID %d", favoriteSongs[0].SongID),
		}

		if err := configs.Database().Create(&recommendation).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add recommendation"})
		}

		return c.Status(fiber.StatusCreated).JSON(recommendation)
	}

	var randomSong models.FavoriteSong
	if err := configs.Database().Where("genre_id = ?", mostLikedGenreID).Order("RANDOM()").First(&randomSong).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to select random song"})
	}

	recommendation := models.Recomendation{
		SongID:     randomSong.SongID, 
		UserID:     userID,            
		GenreID:    mostLikedGenreID,
		Created_at: time.Now(),
		Reason:     fmt.Sprintf("Recommended based on your favorite genre with ID %d", mostLikedGenreID),
	}

	if err := configs.Database().Create(&recommendation).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add recommendation"})
	}

	err = configs.Publish(recommendation.UserID, recommendation.SongID, recommendation.GenreID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to publish recommendation"})
	}

	response := res.RecomendationResponse{
		ID:         recommendation.ID,
		SongID:     recommendation.SongID,
		UserID:     recommendation.UserID,
		GenreID:    recommendation.GenreID,
		Created_at: recommendation.Created_at,
		Reason:     recommendation.Reason,
	}

	return c.JSON(fiber.Map{
		"message": "Recommendation added successfully",
		"data":    response,
	})
}
