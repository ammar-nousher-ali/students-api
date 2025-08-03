package storage

import (
	"github/com/ammar-nousher-ali/students-api/internal/model"
)

type Storage interface {
	//students
	CreateStudent(student model.Student) (int64, error)
	GetStudentById(id int64) (model.Student, error)
	GetStudents() ([]model.Student, error)
	DeleteStudentById(id int64) (int64, error)
	UpdateStudentById(id int64, req model.StudentUpdateRequest) (int64, error)
	SearchStudent(query string) (*[]model.Student, error)

	//users
	CreateUser(user model.User) (int64, error)
	IsEmailTaken(email string) (bool, error)
	GetUserByEmail(email string) (*model.User, error)

	//courses
	CreateCourse(course model.Course) (int64, error)
	GetCourseById(id int64) (*model.Course, error)
	GetAllCourses() ([]model.Course, error)
	UpdateCourse(id int64, req model.CourseUpdateRequest) (*model.Course, error)
}
