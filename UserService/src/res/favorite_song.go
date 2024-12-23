package res

import "time"

type FavoriteSongRes struct {
	ID       int       `json:"id"`
	UserID   int       `json:"user_id"`
	SongID   int       `json:"song_id"`
	GenreID  int       `json:"genre_id"`
	Added_at time.Time `json:"add_at"`
}