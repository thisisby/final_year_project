package data_transfers

import "time"

type CreateNutritionsRequest struct {
	Name    string `json:"name" validate:"required"`
	Value   string `json:"value" validate:"required"`
	OwnerID int    `json:"-"`
}

type UpdateNutritionsRequest struct {
	Name      *string    `json:"name" validate:"omitempty"`
	Value     *string    `json:"value" validate:"omitempty"`
	CreatedAt *time.Time `json:"created_at" validate:"omitempty"`
}

type NutritionsResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Value     string    `json:"value"`
	OwnerID   int       `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
}
