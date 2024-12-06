package data_transfers

import (
	"errors"
	"strings"
)

type CreateWorkoutRequest struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	ExerciseNames []string `json:"exercise_names"`
	OwnerID       int      `json:"-"`
}

type WorkoutsResponse struct {
	ID          int                        `json:"id"`
	Name        string                     `json:"name"`
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
