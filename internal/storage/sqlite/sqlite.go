package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github/com/ammar-nousher-ali/students-api/internal/config"
	"github/com/ammar-nousher-ali/students-api/internal/model"
	"log/slog"
	"strings"

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
	age INTEGER

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

	return &Sqlite{
		Db: db,
	}, nil

}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {

	exists, err := s.checkEmailExists(email)
	if err != nil {
		return 0, err
	}

	if exists {

		return 0, fmt.Errorf("student with this email %s already exists", email)
	}

	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err

	}

	return lastId, nil
}

func (s *Sqlite) checkEmailExists(email string) (bool, error) {
	var count int
	err := s.Db.QueryRow("SELECT COUNT(*) FROM students WHERE email = ?", email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *Sqlite) GetStudentById(id int64) (model.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students WHERE id=? LIMIT 1")

	if err != nil {
		return model.Student{}, err

	}

	defer stmt.Close()

	var student model.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Student{}, fmt.Errorf("no student found for the id %s", fmt.Sprint(id))

		}
		return model.Student{}, fmt.Errorf("query error: %w", err)
	}

	return student, nil

}

func (s *Sqlite) GetStudents() ([]model.Student, error) {
	stmt, err := s.Db.Prepare("SELECT  id, name, email, age FROM students")
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

		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err

		}

		students = append(students, student)
	}

	return students, nil
}

func (s *Sqlite) DeleteStudentById(studentId int64) (int64, error) {

	res, err := s.Db.Exec("DELETE FROM students WHERE id = ?", studentId)

	if err != nil {
		slog.Error("Error deleting student", "err", err)
		return 0, err

	}
	rows, err := res.RowsAffected()
	if err != nil {
		slog.Error("Error checking rows affected", "err", err)
		return 0, err

	}

	if rows == 0 {
		slog.Info("0 Rows")
		return 0, sql.ErrNoRows

	}

	return studentId, nil

}

func (s *Sqlite) UpdateStudentById(studentId int64, req model.StudentUpdateRequest) (int64, error) {
	var fields []string
	var args []any

	if req.Name != nil {
		fields = append(fields, "name = ?")
		args = append(args, *req.Name)
	}

	if req.Email != nil {
		fields = append(fields, "email = ?")
		args = append(args, *&req.Email)

	}

	if req.Age != nil {
		fields = append(fields, "age = ?")
		args = append(args, *&req.Age)

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

func (s *Sqlite) SearchStudent(queryStr string) ([]model.Student, error) {

	dbQuery := "SELECT * FROM students WHERE name LIKE ?"
	rows, err := s.Db.Query(dbQuery, "%"+queryStr+"%")
	if err != nil {
		return []model.Student{}, err
	}

	defer rows.Close()

	var students []model.Student

	for rows.Next() {
		var student model.Student
		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return []model.Student{}, err

		}
		students = append(students, student)
	}

	return students, nil

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
