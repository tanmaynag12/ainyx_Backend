package models

type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=2"`
	Dob  string `json:"dob" validate:"required"`
}

type UpdateUserRequest struct {
	Name string `json:"name" validate:"required,min=2"`
	Dob  string `json:"dob" validate:"required"`
}

type UserResponse struct {
	ID            int32  `json:"id"`
	Name          string `json:"name"`
	Dob           string `json:"dob"`
	Age           *int   `json:"age,omitempty"`
	NextBirthday  *int   `json:"next_birthday_in_days,omitempty"`
}