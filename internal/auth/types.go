package auth

type registerRequest struct {
	Email string `json:"email"`
}

type registerResponse struct {
	APIKey string `json:"apiKey"`
}
