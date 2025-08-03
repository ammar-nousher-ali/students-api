package model

import "time"

type Student struct {
	Id             int64      `json:"id"`
	Name           string     `json:"name" validate:"required"`
	Email          string     `json:"email" validate:"required,email"`
	Age            int        `json:"age" validate:"required"`
	Phone          string     `json:"phone,omitempty"`
	Address        string     `json:"address,omitempty"`
	Gender         string     `json:"gender,omitempty"`
	EnrollmentDate time.Time  `json:"enrollment_date,omitempty"`
	Status         string     `json:"status,omitempty"`
	DeleteAt       *time.Time `json:"delete_at,omitempty"`
}

type StudentUpdateRequest struct {
	//All fields are pointers so you can differentiate between:
	//
	//Field not provided in JSON request (nil)
	//
	//Field provided with a zero value (e.g., "", 0).
	//
	//This is crucial for partial updates.
	Name           *string    `json:"name"`
	Email          *string    `json:"email"`
	Age            *int       `json:"age"`
	Phone          *string    `json:"phone"`
	Address        *string    `json:"address"`
	Gender         *string    `json:"gender"`
	EnrollmentDate *time.Time `json:"enrollment_date"`
	Status         *string    `json:"status"`
}

type CourseUpdateRequest struct {
	CourseCode   *string    `json:"course_code" validate:"required"`
	CourseName   *string    `json:"course_name" validate:"required"`
	Description  *string    `json:"description,omitempty"`
	Credits      *int       `json:"credits" validate:"required"`
	Instructor   *string    `json:"instructor,omitempty"`
	Department   *string    `json:"department,omitempty"`
	Semester     *string    `json:"semester,omitempty"`
	AcademicYear *string    `json:"academic_year,omitempty"`
	Capacity     *int       `json:"capacity,omitempty"`
	Status       *string    `json:"status" validate:"omitempty,oneof=active inactive"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
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

type Course struct {
	Id               int64     `json:"id"`
	CourseCode       string    `json:"course_code" validate:"required"` // unique, required
	CourseName       string    `json:"course_name" validate:"required"` // required
	Description      string    `json:"description,omitempty"`
	Credits          int       `json:"credits" validate:"required"` // required
	Instructor       string    `json:"instructor,omitempty"`
	Department       string    `json:"department,omitempty"`
	Semester         string    `json:"semester,omitempty"`
	AcademicYear     string    `json:"academic_year,omitempty"`
	Capacity         int       `json:"capacity,omitempty"`
	EnrolledStudents []int64   `json:"enrolled_students,omitempty"`
	Status           string    `json:"status" validate:"omitempty,oneof=active inactive"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
	UpdatedAt        time.Time `json:"updated_at,omitempty"`
}
