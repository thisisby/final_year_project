package records

type ExerciseSets struct {
	Record
	WorkoutExercise   WorkoutExercises `db:"workout_exercise"`
	Reps              int              `db:"reps"`
	Weight            float32          `db:"weight"`
	Notes             string           `db:"notes"`
	WorkoutExerciseID int              `db:"workout_exercise_id"`
	OwnerID           int              `db:"owner_id"`
}

type ExerciseDetails struct {
	ExerciseName string `db:"exercise_name"`
	Reps         int    `db:"reps"`
	Weight       int    `db:"weight"`
}
