package student

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github/com/ammar-nousher-ali/students-api/internal/storage"
	"github/com/ammar-nousher-ali/students-api/internal/types"
	"github/com/ammar-nousher-ali/students-api/internal/utils/response"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating a student")

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student) //this is getting info from request and decode it as Student struct
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return

		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return

		}

		// request validation
		if err := validator.New().Struct(student); err != nil {

			validateErrs := err.(validator.ValidationErrors) //err.() this is type assertion in go
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		slog.Info("User created successfully", slog.String("userId", fmt.Sprint(lastId)))

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})

	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.PathValue("id")
		slog.Info("getting a student", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return

		}

		student, err := storage.GetStudentById(intId)
		if err != nil {
			slog.Error("error getting user", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return

		}

		response.WriteJson(w, http.StatusOK, student)

	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("getting all students")
		students, err := storage.GetStudents()
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJson(w, http.StatusOK, students)

	}
}

func DeleteStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("deleting student")
		id := r.PathValue("id")
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return

		}

		slog.Info(fmt.Sprintf("id to be deleted %d", intId))
		deletedStudentId, err := storage.DeleteStudentById(intId)
		if err != nil {
			slog.Info("error while deleting student")

			if errors.Is(err, sql.ErrNoRows) {
				noStudentFoundErr := fmt.Errorf("no student found for the id %d", intId)
				response.WriteJson(w, http.StatusNotFound, response.GeneralError(noStudentFoundErr))
				return

			}

			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		slog.Info(fmt.Sprintf("deleted student id %d", deletedStudentId))
		response.WriteJson(w, http.StatusOK, map[string]any{"message": "Student deleted successfully", "id": deletedStudentId})

	}
}

func UpdateStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("Updating student")

		var req types.StudentUpdateRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			if errors.Is(err, io.EOF) {
				response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
				return
			}

			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("error while parsing json")))
			return
		}

		id := r.PathValue("id")
		studentId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		updatedId, err := storage.UpdateStudentById(studentId, req)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return

		}

		response.WriteJson(w, http.StatusOK, map[string]any{
			"message": "success",
			"id":      updatedId,
		})

	}
}
