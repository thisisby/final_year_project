package records

type Workouts struct {
	Record
	Title       string `db:"title"`
	Description string `db:"description"`
	OwnerID     int    `db:"owner_id"`
	Exercises   []WorkoutExercises
}
