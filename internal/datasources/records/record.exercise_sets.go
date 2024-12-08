package records

type ExerciseSets struct {
	Record
	Reps              int     `db:"reps"`
	Weight            float32 `db:"weight"`
	Notes             string  `db:"notes"`
	WorkoutExerciseID int     `db:"workout_exercise_id"`
	OwnerID           int     `db:"owner_id"`
}
