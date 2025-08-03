package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github/com/ammar-nousher-ali/students-api/internal/config"
	"github/com/ammar-nousher-ali/students-api/internal/model"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3" //here we use under score because we just need this driver we are not using it in code we just need driver repo link
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {

	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER,
	phone TEXT,
	address TEXT,
	gender TEXT,
	enrollment_date TIMESTAMP,
	status TEXT,
	deleted_at TIMESTAMP

	)`)

	if err != nil {
		return nil, err

	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users(
 	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('student', 'teacher'))
	)`)

	if err != nil {
		return nil, err

	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS courses(
 	id INTEGER PRIMARY KEY AUTOINCREMENT,
	course_code TEXT NOT NULL UNIQUE,
    course_name TEXT NOT NULL UNIQUE,
    description TEXT,
    credits INTEGER NOT NULL,
    instructor TEXT,
    department TEXT,
    semester TEXT,
    academic_year TEXT,
    capacity INTEGER,
    status TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
	)`)

	if err != nil {
		return nil, err

	}

	return &Sqlite{
		Db: db,
	}, nil

}

func (s *Sqlite) CreateStudent(student model.Student) (int64, error) {

	exists, err := s.checkEmailExists(student.Email)
	if err != nil {
		return 0, err
	}

	if exists {

		return 0, fmt.Errorf("student with this email %s already exists", student.Email)
	}

	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age, phone, address, gender, enrollment_date, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(student.Name, student.Email, student.Age, student.Phone, student.Address, student.Gender, time.Now(), "active")
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err

	}

	return lastId, nil
}

func (s *Sqlite) GetStudentById(id int64) (model.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age, phone, address, gender, enrollment_date, status FROM students WHERE id=? AND deleted_at IS NULL LIMIT 1")

	if err != nil {
		return model.Student{}, err

	}

	defer stmt.Close()

	var student model.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age, &student.Phone, &student.Address, &student.Gender, &student.EnrollmentDate, &student.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Student{}, err
		}
		return model.Student{}, fmt.Errorf("query error: %w", err)
	}

	return student, nil

}

func (s *Sqlite) GetStudents() ([]model.Student, error) {
	stmt, err := s.Db.Prepare("SELECT  id, name, email, age, phone, address, gender, enrollment_date, status FROM students WHERE deleted_at IS NULL")
	if err != nil {
		return nil, err

	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err

	}

	defer rows.Close()

	var students []model.Student

	for rows.Next() {
		var student model.Student

		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age, &student.Phone, &student.Address, &student.Gender, &student.EnrollmentDate, &student.Status)
		if err != nil {
			return nil, err

		}

		students = append(students, student)
	}

	return students, nil
}

func (s *Sqlite) DeleteStudentById(studentId int64) (int64, error) {

	//res, err := s.Db.Exec("DELETE FROM students WHERE id = ?", studentId)
	res, err := s.Db.Exec("UPDATE students SET deleted_at = ? WHERE id = ?", time.Now(), studentId)

	if err != nil {
		return 0, err

	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err

	}

	if rows == 0 {
		return 0, sql.ErrNoRows

	}

	return studentId, nil

}

func (s *Sqlite) UpdateStudentById(studentId int64, req model.StudentUpdateRequest) (int64, error) {
	var fields []string
	var args []any

	//Why *req.Name and not req.Name?
	//Because req.Name is of type *string (a pointer), and you need the actual value (type string) to pass as a query argument.
	if req.Name != nil {
		fields = append(fields, "name = ?")
		args = append(args, *req.Name)
	}

	if req.Email != nil {
		fields = append(fields, "email = ?")
		args = append(args, *req.Email)

	}

	if req.Age != nil {
		fields = append(fields, "age = ?")
		args = append(args, *req.Age)

	}

	if req.Phone != nil {
		fields = append(fields, "phone = ?")
		args = append(args, *req.Phone)
	}

	if req.Address != nil {
		fields = append(fields, "address = ?")
		args = append(args, *req.Address)
	}

	if req.Gender != nil {
		fields = append(fields, "gender = ?")
		args = append(args, *req.Gender)
	}

	if req.EnrollmentDate != nil {
		fields = append(fields, "enrollment_date = ?")
		args = append(args, *req.EnrollmentDate)
	}

	if req.Status != nil {
		fields = append(fields, "status =?")
		args = append(args, *req.Status)
	}

	if len(fields) == 0 {
		return 0, fmt.Errorf("no fields to update")

	}

	args = append(args, studentId)
	query := fmt.Sprintf("UPDATE students SET %s WHERE id = ?", strings.Join(fields, ", "))

	result, err := s.Db.Exec(query, args...)
	if err != nil {
		return 0, err

	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err

	}

	if rows == 0 {

		return 0, sql.ErrNoRows

	}

	return studentId, nil

}

func (s *Sqlite) SearchStudent(queryStr string) (*[]model.Student, error) {

	dbQuery := "SELECT * FROM students WHERE name LIKE ?"
	rows, err := s.Db.Query(dbQuery, "%"+queryStr+"%")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []model.Student

	for rows.Next() {
		var student model.Student
		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age, &student.Phone, &student.Address, &student.Gender, &student.EnrollmentDate, &student.Status)
		if err != nil {
			return nil, err

		}
		students = append(students, student)
	}

	if len(students) == 0 {
		return nil, sql.ErrNoRows

	}

	return &students, nil

}

func (s *Sqlite) IsEmailTaken(email string) (bool, error) {

	var count int
	row := s.Db.QueryRow("SELECT COUNT(*) from users WHERE email = ?", email)
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *Sqlite) CreateUser(user model.User) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, ?)")

	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(user.Name, user.Email, user.Password, user.Role)
	if err != nil {
		return 0, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err

	}

	return lastId, nil
}

func (s *Sqlite) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	row := s.Db.QueryRow("SELECT id, name, email, password, role from users WHERE email = ? LIMIT 1", email)

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.New("user not found with this email")
			return nil, err
		}

		return nil, err
	}

	return &user, nil

}

func (s *Sqlite) checkEmailExists(email string) (bool, error) {
	var count int
	err := s.Db.QueryRow("SELECT COUNT(*) FROM students WHERE email = ?", email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

//course

func (s *Sqlite) CreateCourse(course model.Course) (int64, error) {

	//slog.Info("course", "struct", course)
	result, err := s.Db.Exec("INSERT INTO courses (course_code, course_name, description, credits, instructor, department, semester, academic_year, capacity, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		course.CourseCode, course.CourseName, course.Description, course.Credits,
		course.Instructor, course.Department, course.Semester, course.AcademicYear,
		course.Capacity, course.Status, course.CreatedAt, course.UpdatedAt)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil

}

func (s *Sqlite) GetCourseById(id int64) (*model.Course, error) {

	var course model.Course

	row := s.Db.QueryRow("SELECT id, course_code, course_name, description, credits, instructor, department, semester, academic_year, capacity, status, created_at, updated_at from courses WHERE id = ?", id)

	err := row.Scan(&course.Id, &course.CourseCode, &course.CourseName, &course.Description, &course.Credits,
		&course.Instructor, &course.Department, &course.Semester, &course.AcademicYear,
		&course.Capacity, &course.Status, &course.CreatedAt, &course.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &course, nil

}

func (s *Sqlite) GetAllCourses() ([]model.Course, error) {

	var courses []model.Course

	rows, err := s.Db.Query("SELECT * from courses")
	if err != nil {
		return courses, err
	}
	//rows.Next() moves to the next row in the result set.
	//It returns true as long as more rows are available.
	for rows.Next() {

		var course model.Course

		err := rows.Scan(&course.Id, &course.CourseCode, &course.CourseName, &course.Description, &course.Credits,
			&course.Instructor, &course.Department, &course.Semester, &course.AcademicYear,
			&course.Capacity, &course.Status, &course.CreatedAt, &course.UpdatedAt)

		if err != nil {
			return nil, err
		}

		courses = append(courses, course)

	}

	return courses, nil

}

func (s *Sqlite) UpdateCourse(id int64, req model.CourseUpdateRequest) (*model.Course, error) {
	var fields []string
	var args []any

	if req.CourseCode != nil {
		fields = append(fields, "course_code = ?")
		args = append(args, *req.CourseCode)
	}

	if req.CourseName != nil {
		fields = append(fields, "course_name = ?")
		args = append(args, *req.CourseName)
	}

	if req.Description != nil {
		fields = append(fields, "description = ?")
		args = append(args, *req.Description)
	}

	if req.Credits != nil {
		fields = append(fields, "credits = ?")
		args = append(args, *req.Credits)
	}

	if req.Instructor != nil {
		fields = append(fields, "instructor = ?")
		args = append(args, *req.Instructor)
	}

	if req.Department != nil {
		fields = append(fields, "department = ?")
		args = append(args, *req.Department)
	}

	if req.Semester != nil {
		fields = append(fields, "semester = ?")
		args = append(args, *req.Semester)
	}

	if req.AcademicYear != nil {
		fields = append(fields, "academic_year = ?")
		args = append(args, *req.AcademicYear)
	}

	if req.Capacity != nil {
		fields = append(fields, "capacity = ?")
		args = append(args, *req.Capacity)
	}

	if req.Status != nil {
		fields = append(fields, "status = ?")
		args = append(args, *req.Status)
	}

	if req.UpdatedAt != nil {
		fields = append(fields, "updated_at = ?")
		args = append(args, *req.UpdatedAt)
	}
	args = append(args, id)
	query := fmt.Sprintf("UPDATE courses SET %s WHERE id = ?", strings.Join(fields, ", "))

	_, err := s.Db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	course, err := s.GetCourseById(id)
	if err != nil {
		return nil, err
	}

	return course, nil

}

func (s *Sqlite) DeleteCourseById(id int64) (int64, error) {

	res, err := s.Db.Exec("DELETE from courses WHERE id = ?", id)
	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	if rows == 0 {

		return 0, sql.ErrNoRows
	}

	return id, nil

}
