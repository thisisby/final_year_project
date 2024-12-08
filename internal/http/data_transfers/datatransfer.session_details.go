package data_transfers

type CreateSessionDetailsRequest struct {
	SessionID int    `json:"session_id" validate:"required"`
	Name      string `json:"name" validate:"required"`
	Value     string `json:"value" validate:"required"`
}

type CreateSessionDetailsResponse struct {
	SessionID int    `json:"session_id"`
	Name      string `json:"name"`
	Value     string `json:"value"`
}

type UpdateSessionDetailsRequest struct {
	Name  *string `json:"name" validate:"omitempty"`
	Value *string `json:"value" validate:"omitempty"`
}
