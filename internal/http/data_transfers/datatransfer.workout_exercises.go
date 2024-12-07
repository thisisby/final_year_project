package data_transfers

type WorkoutExercisesResponse struct {
	Exercise      ExercisesResponse `json:"exercise"`
	ID            int               `json:"id"`
	MainNote      string            `json:"main_note"`
	SecondaryNote string            `json:"secondary_note"`
	WorkoutID     int               `json:"workout_id"`
	OwnerID       int               `json:"owner_id"`
	ExerciseID    int               `json:"exercise_id"`
}

type CreateWorkoutExercisesRequest struct {
	MainNote      string `json:"main_note" validate:"omitempty"`
	SecondaryNote string `json:"secondary_note" validate:"omitempty"`
	WorkoutID     int    `json:"workout_id" validate:"required"`
	OwnerID       int    `json:"-"`
	ExerciseID    int    `json:"exercise_id" validate:"required"`
}

type UpdateWorkoutExercisesRequest struct {
	MainNote      *string `json:"main_note" validate:"omitempty"`
	SecondaryNote *string `json:"secondary_note" validate:"omitempty"`
}
