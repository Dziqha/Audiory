package res


type AdminResponse struct {
	ID int `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
}
type AdminResponseUpdate struct {
	ID int `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type AdminResponseLogin struct {
	ID int `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`	
	Token string `json:"token"`	
}