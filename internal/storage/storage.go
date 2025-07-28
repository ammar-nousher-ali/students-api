package storage

import (
	"github/com/ammar-nousher-ali/students-api/internal/model"
)

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (model.Student, error)
	GetStudents() ([]model.Student, error)
	DeleteStudentById(id int64) (int64, error)
	UpdateStudentById(id int64, req model.StudentUpdateRequest) (int64, error)
	SearchStudent(query string) ([]model.Student, error)
	IsEmailTaken(email string) (bool, error)
	CreateUser(user model.User) (int64, error)
}
