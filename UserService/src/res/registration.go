package res

import (
	"time"
)

type RegistrationResponse struct {
	ID        int       `json:"id"`
	ArtistID  int       `json:"artist_id"`
	AdminID   int       `json:"admin_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}