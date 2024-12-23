package controller

import (
	"UserService/configs"
	"UserService/proto"
	"UserService/src/models"
	"UserService/src/res"
	"UserService/src/utils"
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

type RegistrationController struct{}

func NewRegistrationController() *RegistrationController {
	return &RegistrationController{}
}

func (r *RegistrationController) AddRegister(c *fiber.Ctx) error {
	var reqData models.RegistrationRequest
	var admin models.Admin

	if err := c.BodyParser(&reqData); err != nil {
		return errors.Wrap(err, "failed to parse request body")
	}

	if err := utils.ValidateStruct(&reqData); err != nil {
		return errors.Wrap(err, "failed to validate registration")
	}

	if err := configs.Database().First(&admin, reqData.AdminID).Error; err != nil {
		return errors.New("admin not found")
	}

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to connect to gRPC server",
		})
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing gRPC connection: %v", err)
		}
	}()
	log.Println("Connected to gRPC server")

	artistClient := proto.NewArtistServiceClient(conn)
	log.Println("artist client initialized")

	if reqData.ArtistID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Artist ID is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	artistResponse, err := artistClient.GetArtistById(ctx, &proto.ArtistRequest{Id: int32(reqData.ArtistID)})
	if err != nil {
		log.Printf("Error calling gRPC service: %v", err) 
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get song from gRPC service",
			"error":   err.Error(), 
		})
	}


	req := models.Registration{
		ArtistID:  int(artistResponse.Id),
		AdminID:   reqData.AdminID,
		Status:    reqData.Status,
		CreatedAt: time.Now(),
	}
	if err := configs.Database().Create(&req).Error; err != nil {
		return errors.Wrap(err, "failed to create registration")
	}
	res := res.RegistrationResponse{
		ID:        req.ID,
		ArtistID:  req.ArtistID,
		AdminID:   req.AdminID,
		Status:    req.Status,
		CreatedAt: time.Now(),
	}

	return c.JSON(fiber.Map{
		"message": "Successfully added registration",
		"data":    res,
	})
}

func (r *RegistrationController)  GetAllRegistration(c *fiber.Ctx) error {
	var data []models.Registration
	var resData []res.RegistrationResponse

	err := configs.Database().Transaction(func(tx *gorm.DB) error {
		if result := tx.Find(&data); result.Error != nil {
			return errors.Wrap(result.Error, "failed to get all registration")
		}

		return nil
	})

	if len(data) == 0 {
		return errors.Wrap(err, "favorite song not found")
	}

	if err != nil {
		return errors.Wrap(err, "failed to process the request")
	} 

	for _, findata := range data {
		resData = append(resData, res.RegistrationResponse{
			ID:        findata.ID,
			ArtistID:  findata.ArtistID,
			AdminID:   findata.AdminID,
			Status:    findata.Status,
			CreatedAt: findata.CreatedAt,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Successfully get all registration",
		"data":    resData,
	})
}


func (r *RegistrationController)  GetRegistrationById(c *fiber.Ctx) error {
	var data models.Registration
	var resData res.RegistrationResponse
	id := c.Params("id")

	err := configs.Database().Transaction(func(tx *gorm.DB) error {
		if result := tx.Where("id = ?", id).First(&data); result.Error != nil {
			return errors.Wrap(result.Error, "failed to get registration by id")
		}
		return nil
	})

	if err != nil {
		return errors.Wrap(err, "failed to process the request")
	}

	resData = res.RegistrationResponse{
		ID:        data.ID,
		ArtistID:  data.ArtistID,
		AdminID:   data.AdminID,
		Status:    data.Status,
		CreatedAt: data.CreatedAt,
	}

	return c.JSON(fiber.Map{
		"message": "Successfully get registration by id",
		"data":    resData,
	})
}


func (r *RegistrationController) UpdateRegistration(c *fiber.Ctx) error {
	var data models.Registration
	var req models.RegistrationRequest
	id := c.Params("id")

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	if result := configs.Database().Where("id = ?", id).First(&data); result.Error != nil {
		return errors.Wrap(result.Error, "failed to get registration by id")
	}

	if req.AdminID != 0 {
		admin := configs.Database().Where("id = ?", req.AdminID).First(&models.Admin{})
		if admin.Error != nil {
			return errors.Wrap(admin.Error, "failed to get admin by id")
		}
		data.AdminID = req.AdminID
	}

	if req.ArtistID != 0 {
		conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to connect to gRPC server for artist service",
			})
		}
		defer func() {
			if err := conn.Close(); err != nil {
				log.Printf("Error closing gRPC connection: %v", err)
			}
		}()
		log.Println("Connected to gRPC server for artist service")

		artistClient := proto.NewArtistServiceClient(conn)
		log.Println("Artist client initialized")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		artistResponse, err := artistClient.GetArtistById(ctx, &proto.ArtistRequest{Id: int32(req.ArtistID)})
		if err != nil {
			log.Printf("Error calling gRPC service for artist: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to get artist from gRPC service",
				"error":   err.Error(),
			})
		}
		if artistResponse.Id == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Artist not found",
			})
		}
		data.ArtistID = int(artistResponse.Id)
	}

	log.Println("Artist ID:", data.ArtistID)

	if req.Status != "" {
		data.Status = req.Status
	}

	err := configs.Database().Transaction(func(tx *gorm.DB) error {
		if result := tx.Save(&data); result.Error != nil {
			return errors.Wrap(result.Error, "failed to update registration")
		}
		return nil
	})

	if err != nil {
		return errors.Wrap(err, "failed to process the request")
	}

	responseData := res.RegistrationResponse{
		ID:        data.ID,
		ArtistID:  data.ArtistID,
		AdminID:   data.AdminID,
		Status:    data.Status,
		CreatedAt: data.CreatedAt,
	}

	return c.JSON(fiber.Map{
		"message": "Registration updated successfully",
		"data":    responseData,
	})
}


func (r *RegistrationController)  DeleteRegistration(c *fiber.Ctx) error {
	var data models.Registration
	id := c.Params("id")

	err := configs.Database().Transaction(func(tx *gorm.DB) error {
		if result := tx.Where("id = ?", id).First(&data); result.Error != nil {
			return errors.Wrap(result.Error, "failed to get registration by id")
		}

		if result := tx.Delete(&data); result.Error != nil {
			return errors.Wrap(result.Error, "failed to delete registration")
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "failed to process the request")
	}

	return c.JSON(fiber.Map{
		"message": "Registration deleted successfully",
	})
}
