package data_transfers

import (
	"errors"
	"strings"
)

type CreateWorkoutRequest struct {
	Title         string   `json:"title" validate:"required"`
	Description   string   `json:"description" validate:"omitempty"`
	ExerciseNames []string `json:"exercise_names" validate:"required"`
	OwnerID       int      `json:"-"`
}

type UpdateWorkoutRequest struct {
	Title       *string `json:"title" validate:"omitempty"`
	Description *string `json:"description" validate:"omitempty"`
}

type AddExercisesRequest struct {
	ExerciseNames []string `json:"exercise_names"`
}

type DeleteExerciseRequest struct {
	ExerciseIDs []int `json:"exercise_ids"`
}

type WorkoutsResponse struct {
	ID          int                        `json:"id"`
	Title       string                     `json:"title"`
	Description string                     `json:"description"`
	OwnerID     int                        `json:"owner_id"`
	Exercises   []WorkoutExercisesResponse `json:"exercises"`
}

func (r *CreateWorkoutRequest) Validate() error {
	if len(r.ExerciseNames) > 0 {
		for _, exercise := range r.ExerciseNames {
			if strings.TrimSpace(exercise) == "" {
				return errors.New("exercise_names contains an empty string")
			}
		}
	}
	return nil
}
