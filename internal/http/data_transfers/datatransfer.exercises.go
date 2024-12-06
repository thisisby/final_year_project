package data_transfers

type CreateExercisesRequest struct {
	Name string `json:"name" validate:"required"`
}

type ExercisesResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UpdateExercisesRequest struct {
	Name *string `json:"name" validate:"omitempty"`
}
