package data_transfers

import "time"

type CreateExerciseSetsRequest struct {
	Reps              int     `json:"reps"`
	Weight            float32 `json:"weight"`
	WorkoutExerciseID int     `json:"workout_exercise_id"`
	OwnerID           int     `json:"-"`
}

type ExerciseSetsResponse struct {
	ID                int       `json:"id"`
	Reps              int       `json:"reps"`
	Weight            float32   `json:"weight"`
	Notes             string    `json:"notes"`
	WorkoutExerciseID int       `json:"workout_exercise_id"`
	OwnerID           int       `json:"owner_id"`
	CreatedAt         time.Time `json:"created_at"`
}

type UpdateExerciseSetsRequest struct {
	Reps      *int     `json:"reps"`
	Weight    *float32 `json:"weight"`
	Notes     *string  `json:"notes"`
	CreatedAt *string  `json:"created_at"`
}
