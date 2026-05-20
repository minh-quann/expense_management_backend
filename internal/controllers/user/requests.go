package user

type UpdateProfileRequest struct {
	DisplayName  string `json:"display_name" binding:"required"`
	CurrencyCode string `json:"currency_code" binding:"required"`
	PhoneNumber  string `json:"phone_number"`
	Address      string `json:"address"`
	Gender       string `json:"gender"`
}
