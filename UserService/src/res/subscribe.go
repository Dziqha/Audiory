package res

import "time"

type SubscriptionResponse struct {
	SubscriptionStart time.Time  `json:"subscription_start"`
	SubscriptionEnd   time.Time  `json:"subscription_end"`
	SubscriptionToken string     `json:"subscription_token"`
	SubscriptionType  string     `json:"subscription_type"`
	IsActive          bool       `json:"is_active"`
	UserID            int        `json:"user_id"`
	PaymentToken      string     `json:"payment_token"`
	PaymentRedirectURL string     `json:"payment_redirect_url"`}
