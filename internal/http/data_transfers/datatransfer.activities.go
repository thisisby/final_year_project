package data_transfers

type CreateActivityRequest struct {
	Name            string `json:"name" validate:"required"`
	ActivityGroupID int    `json:"activity_group_id" validate:"required"`
}

type UpdateActivityRequest struct {
	Name            *string `json:"name" validate:"omitempty"`
	ActivityGroupID *int    `json:"activity_group_id" validate:"omitempty"`
}

type ActivityResponse struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	ActivityGroupID int    `json:"activity_group_id"`
}
