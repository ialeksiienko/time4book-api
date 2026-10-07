package handlers

type ErrorResponse struct {
	Status bool   `json:"status" example:"false"`
	Error  string `json:"error" example:"error message detail"`
}

type SuccessResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message,omitempty" example:"success"`
}

