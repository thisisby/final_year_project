package records

type WorkoutExercises struct {
	Record
	Exercise      Exercises `db:"exercise"`
	MainNote      string    `db:"main_note"`
	SecondaryNote string    `db:"secondary_note"`
	WorkoutID     int       `db:"workout_id"`
	OwnerID       int       `db:"owner_id"`
	ExerciseID    int       `db:"exercise_id"`
}
