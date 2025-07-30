package model

import "time"

type Student struct {
	Id             int64     `json:"id"`
	Name           string    `json:"name" validate:"required"`
	Email          string    `json:"email" validate:"required,email"`
	Age            int       `json:"age" validate:"required"`
	Phone          string    `json:"phone,omitempty"`
	Address        string    `json:"address,omitempty"`
	Gender         string    `json:"gender,omitempty"`
	EnrollmentDate time.Time `json:"enrollment_date,omitempty"`
	Status         string    `json:"status,omitempty"`
}

type StudentUpdateRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
	Age   *int    `json:"age"`
}

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role"`
}
type Creds struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
