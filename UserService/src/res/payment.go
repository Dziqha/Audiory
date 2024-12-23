package res

type PaymentResponse struct {
	Token string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}