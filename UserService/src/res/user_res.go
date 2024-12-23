package res

type UserRes struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserLoginRes struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Username string `json:"username"`
	Token string `json:"token"`
}