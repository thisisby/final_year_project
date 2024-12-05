package data_transfers

type CreateUsersRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UsersResponse struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"_"`
}

type UpdateUsersRequest struct {
	Email    *string `json:"email" validate:"omitempty,email"`
	Password *string `json:"password" validate:"omitempty,min=8"`
}
