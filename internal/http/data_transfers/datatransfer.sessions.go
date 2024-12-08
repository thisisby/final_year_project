package data_transfers

import "time"

type CreateSessionRequest struct {
	Notes      string    `json:"notes"`
	StartTime  time.Time `json:"start_time" validate:"required"`
	EndTime    time.Time `json:"end_time" validate:"required"`
	ActivityID int       `json:"activity_id" validate:"required"`
	OwnerID    int       `json:"-"`
}

type UpdateSessionRequest struct {
	Notes      *string    `json:"notes" validate:"omitempty"`
	StartTime  *time.Time `json:"start_time" validate:"omitempty"`
	EndTime    *time.Time `json:"end_time" validate:"omitempty"`
	ActivityID *int       `json:"activity_id" validate:"omitempty"`
}

type SessionResponse struct {
	Activity   ActivityResponse `json:"activity"`
	ID         int              `json:"id"`
	Notes      string           `json:"notes"`
	StartTime  time.Time        `json:"start_time"`
	EndTime    time.Time        `json:"end_time"`
	ActivityID int              `json:"activity_id"`
}
