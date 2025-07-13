package storage

import "github/com/ammar-nousher-ali/students-api/internal/types"




type Storage interface{
	CreateStudent(name string, email string, age int)(int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetStudents()([]types.Student, error)
	DeleteStudentById(id int64) (int64,error)
}