package res

import "time"

type ListeningHistoryRes struct {
	ID             int       `json:"id"`
	PlayedAt       time.Time `json:"played_at"`
	DurationPlayed int       `json:"duration_played"`
	UserID         int       `json:"user_id"`
	SongID         int       `json:"song_id"`
}