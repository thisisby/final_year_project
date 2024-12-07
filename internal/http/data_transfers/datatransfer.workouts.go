package data_transfers

type CreateWorkoutRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"omitempty"`
	OwnerID     int    `json:"-"`
}

type UpdateWorkoutRequest struct {
	Title       *string `json:"title" validate:"omitempty"`
	Description *string `json:"description" validate:"omitempty"`
}

type WorkoutsResponse struct {
	ID          int                        `json:"id"`
	Title       string                     `json:"title"`
	Description string                     `json:"description"`
	OwnerID     int                        `json:"owner_id"`
	Exercises   []WorkoutExercisesResponse `json:"exercises"`
}
