package records

type Workouts struct {
	Record
	Title       string  `db:"title"`
	Description string  `db:"description"`
	IsPrivate   bool    `db:"is_private"`
	Price       float64 `db:"price"`
	OwnerID     int     `db:"owner_id"`
	Exercises   []WorkoutExercises
}
