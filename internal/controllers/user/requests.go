package user

type UpdateProfileRequest struct {
	DisplayName  string `json:"display_name" binding:"required"`
	CurrencyCode string `json:"currency_code" binding:"required"`
	PhoneNumber  string `json:"phone_number"`
	Address      string `json:"address"`
	Gender       string `json:"gender"`
}

type SetPINRequest struct {
	Pin              string `json:"pin" binding:"required,min=4,max=6"`
	SecurityQuestion string `json:"security_question" binding:"required"`
	SecurityAnswer   string `json:"security_answer" binding:"required"`
}

type VerifyPINRequest struct {
	Pin string `json:"pin" binding:"required"`
}

type ResetPINRequest struct {
	SecurityAnswer string `json:"security_answer" binding:"required"`
	NewPin         string `json:"new_pin" binding:"required,min=4,max=6"`
}

type DisablePINRequest struct {
	Pin string `json:"pin" binding:"required"`
}
