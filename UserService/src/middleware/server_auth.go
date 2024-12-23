package middleware

import (
	"context"
	"errors"
	"os"
	"strconv"
	"strings"

	pb "UserService/proto"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/gofiber/fiber/v2"
)

type AuthServiceServer struct {
	pb.UnimplementedAuthServiceServer
}

func (s *AuthServiceServer) ValidateTokenUser(ctx context.Context, req *pb.ValidateTokenRequestUser) (*pb.ValidateTokenResponseUser, error) {
	tokenSecret := os.Getenv("TOKEN_SECRET_USER")
	tokenString := strings.TrimSpace(req.Token)

	if tokenString == "" {
		return &pb.ValidateTokenResponseUser{
			IsValid: false,
			Error:   "Unauthorized: missing token",
		}, nil
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("unauthorized: invalid signing method")
				}
				return []byte(tokenSecret), nil
			})

	if err != nil {
		return &pb.ValidateTokenResponseUser{
			IsValid: false,
			Error:   "Invalid or expired token: " + err.Error(),
		}, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return &pb.ValidateTokenResponseUser{
			IsValid: false,
			Error:   "Invalid token claims",
		}, nil
	}

	userId, ok := claims["user_id"]
	if !ok {
		return &pb.ValidateTokenResponseUser{
			IsValid: false,
			Error:   "Invalid token claims: user_id is not a string",
		}, nil
	}

	return &pb.ValidateTokenResponseUser{
		IsValid: true,
  		UserId: strconv.Itoa(int(userId.(float64))),
	}, nil
}
func (s *AuthServiceServer) ValidateTokenAdmin(ctx context.Context, req *pb.ValidateTokenRequestAdmin) (*pb.ValidateTokenResponseAdmin, error) {
	tokenSecret := os.Getenv("TOKEN_SECRET_ADMIN")
	tokenString := strings.TrimSpace(req.Token)

	if tokenString == "" {
		return &pb.ValidateTokenResponseAdmin{
			IsValid: false,
			Error:   "Unauthorized: missing token",
		}, nil
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("unauthorized: invalid signing method")
				}
				return []byte(tokenSecret), nil
			})

	if err != nil {
		return &pb.ValidateTokenResponseAdmin{
			IsValid: false,
			Error:   "Invalid or expired token: " + err.Error(),
		}, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return &pb.ValidateTokenResponseAdmin{
			IsValid: false,
			Error:   "Invalid token claims",
		}, nil
	}

	adminId, ok := claims["admin_id"]
	if !ok {
		return &pb.ValidateTokenResponseAdmin{
			IsValid: false,
			Error:   "Invalid token claims: admin_id is not a string",
		}, nil
	}

	return &pb.ValidateTokenResponseAdmin{
		IsValid: true,
  		AdminId: strconv.Itoa(int(adminId.(float64))),
	}, nil
}
