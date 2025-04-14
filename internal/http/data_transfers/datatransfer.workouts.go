package data_transfers

type CreateWorkoutRequest struct {
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description" validate:"omitempty"`
	IsPrivate   bool    `json:"is_private" validate:"omitempty"`
	Price       float64 `json:"price" validate:"omitempty"`
	OwnerID     int     `json:"-"`
}

type UpdateWorkoutRequest struct {
	Title       *string  `json:"title" validate:"omitempty"`
	Description *string  `json:"description" validate:"omitempty"`
	IsPrivate   *bool    `json:"is_private" validate:"omitempty"`
	Price       *float64 `json:"price" validate:"omitempty"`
}

type WorkoutsResponse struct {
	ID          int                        `json:"id"`
	Title       string                     `json:"title"`
	Description string                     `json:"description"`
	IsPrivate   bool                       `json:"is_private"`
	Price       float64                    `json:"price"`
	OwnerID     int                        `json:"owner_id"`
	LikesCount  int                        `json:"likes_count"`
	Exercises   []WorkoutExercisesResponse `json:"exercises"`
}

type WorkoutGenerateRequest struct {
	Level     string   `json:"level" validate:"required"`
	Goal      string   `json:"goal" validate:"required"`
	BodyAreas []string `json:"body_areas" validate:"required"`
	Gender    string   `json:"gender" validate:"required"`
	Age       string   `json:"age" validate:"required"`
	Details   string   `json:"details" validate:"omitempty"`
	OwnerID   int      `json:"-"`
}
