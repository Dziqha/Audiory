package res

import "time"

type ReviewResponse struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	AlbumID   int       `json:"album_id"`
	Rating    string    `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}