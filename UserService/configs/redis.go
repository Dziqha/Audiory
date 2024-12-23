package configs

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)


func Initialize() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", 
        os.Getenv("REDIS_HOST"), 
        os.Getenv("REDIS_PORT")),
	})

	if _,err := client.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}else {
		log.Println("Connected to Redis")
	}

	return client
}

func Publish(userId int, songId int, genreId int) error {
	message := fmt.Sprintf("New recommendation for user %d: Song ID %d (Genre ID %d)", userId, songId, genreId)
	err := Initialize().Publish(context.Background(), "recommendation_channel", message).Err()
	if err != nil {
		return errors.Wrap(err, "failed to publish message")
	}

	return nil
}

func Subscribe(){
	client := Initialize().Subscribe(context.Background(), "recommendation_channel")
	defer client.Close()

	for {
		msg, err := client.ReceiveMessage(context.Background())
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			continue
		}
		 log.Printf("New Recomendation Received message: %s", msg.Payload)
	}
}