package records

type WorkoutExercises struct {
	Record
	Name          string `db:"name"`
	MainNote      string `db:"main_note"`
	SecondaryNote string `db:"secondary_note"`
	WorkoutID     int    `db:"workout_id"`
}
