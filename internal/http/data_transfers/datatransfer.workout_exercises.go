package data_transfers

type WorkoutExercisesResponse struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	MainNote      string `json:"main_note"`
	SecondaryNote string `json:"secondary_note"`
	WorkoutID     int    `json:"workout_id"`
}
