package records

type Exercises struct {
	Record
	Name string `db:"name"`
}

// ExercisesWithWorkoutCheck is a struct that is not stored in the database.
// It is used to check if an exercise is in a workout.
type ExercisesWithWorkoutCheck struct {
	Record
	Name        string `db:"name"`
	IsInWorkout bool   `db:"is_in_workout"`
}
