package res

import "time"

type RecomendationResponse struct {
	ID         int       `json:"id"`
	SongID     int       `json:"song_id"`
	UserID     int       `json:"user_id"`
	GenreID    int       `json:"genre_id"`
	Created_at time.Time `json:"created_at"`
	Reason     string    `json:"reason"`
}