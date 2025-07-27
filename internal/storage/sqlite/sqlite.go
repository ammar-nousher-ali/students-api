package sqlite

import (
	"database/sql"
	"fmt"
	"github/com/ammar-nousher-ali/students-api/internal/config"
	"github/com/ammar-nousher-ali/students-api/internal/types"

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

	return &Sqlite{
		Db: db,
	}, nil

}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {

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

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students WHERE id=? LIMIT 1")

	if err != nil {
		return types.Student{}, err

	}

	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found for the id %s", fmt.Sprint(id))

		}
		return types.Student{}, fmt.Errorf("query error: %w", err)
	}

	return student, nil

}

func (s *Sqlite) GetStudents() ([]types.Student, error) {
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

	var students []types.Student

	for rows.Next() {
		var student types.Student

		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err

		}

		students = append(students, student)
	}

	return students, nil
}

func (s *Sqlite) DeleteStudentById(studentId int64) (int64, error) {

	stmt, err := s.Db.Prepare("DELETE FROM students WHERE id=?")

	if err != nil {
		return 0, err

	}

	defer stmt.Close()

	_, err = stmt.Exec(studentId)

	if err != nil {
		return 0, err

	}

	return studentId, nil

}
