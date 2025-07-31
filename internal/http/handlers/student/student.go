package student

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github/com/ammar-nousher-ali/students-api/internal/model"
	"github/com/ammar-nousher-ali/students-api/internal/storage"
	"github/com/ammar-nousher-ali/students-api/internal/utils/response"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating a student")

		var student model.Student

		err := json.NewDecoder(r.Body).Decode(&student) //this is getting info from request and decode it as Student struct
		if errors.Is(err, io.EOF) {
			response.WriteJson(w,
				http.StatusBadRequest,
				response.GeneralError(
					fmt.Errorf("empty body"),
					http.StatusBadRequest,
				),
			)
			return

		}

		if err != nil {
			// response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			response.WriteJson(w,
				http.StatusBadRequest,
				response.GeneralError(
					err,
					http.StatusBadRequest,
				),
			)
			return

		}

		// request validation
		if err := validator.New().Struct(student); err != nil {

			var validateErrs validator.ValidationErrors
			errors.As(err, &validateErrs) //err.() this is type assertion in go
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs, http.StatusBadRequest))
			return
		}

		lastId, err := storage.CreateStudent(student)

		slog.Info("student created successfully", slog.String("userId", fmt.Sprint(lastId)))

		if err != nil {

			if strings.Contains(err.Error(), "already") {
				response.WriteJson(w,
					http.StatusConflict,
					response.GeneralError(
						err,
						http.StatusConflict,
					),
				)
				return
			}

			response.WriteJson(w,
				http.StatusInternalServerError,
				response.GeneralError(
					err,
					http.StatusInternalServerError,
				),
			)
			return
		}

		response.WriteJson(w, http.StatusCreated,
			response.GeneralResponse(
				"Student created successfully",
				http.StatusCreated,
				map[string]int64{"id": lastId},
			),
		)

	}
}

func NewBatch(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var students []model.Student

		err := json.NewDecoder(r.Body).Decode(&students)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request"), http.StatusBadRequest))
			return
		}

		if len(students) == 0 {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("please provide at least one student data to create"), http.StatusBadRequest))
			return
		}

		var batchResponse response.BatchResponse
		for _, student := range students {
			id, err := storage.CreateStudent(student)

			if err != nil {
				batchResponse.Data = append(batchResponse.Data, response.BatchData{
					Success: false,
					Data: map[string]any{
						"message": "failed",
						"reason":  err.Error(),
					},
				})

			} else {
				batchResponse.Data = append(batchResponse.Data, response.BatchData{
					Success: true,
					Data: map[string]any{
						"message": "success",
						"id":      id,
					},
				})
			}

		}

		response.WriteJson(w, http.StatusOK, response.GeneralBatchResponse("success", http.StatusCreated, batchResponse.Data))

	}

}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.PathValue("id")
		slog.Info("getting a student", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			// response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			response.WriteJson(w,
				http.StatusBadRequest,
				response.GeneralError(
					fmt.Errorf("empty body"),
					http.StatusBadRequest,
				),
			)
			return

		}

		student, err := storage.GetStudentById(intId)
		if err != nil {
			slog.Error("error getting user", slog.String("id", id))

			if errors.Is(err, sql.ErrNoRows) {
				slog.Error("sql no rows err", "error", err)
				response.WriteJson(w, http.StatusNotFound, response.GeneralError(fmt.Errorf("no student found for the given id = %s", id), http.StatusNotFound))
				return
			}

			response.WriteJson(w,
				http.StatusInternalServerError,
				response.GeneralError(
					fmt.Errorf("empty body"),
					http.StatusInternalServerError,
				),
			)
			return

		}

		response.WriteJson(w, http.StatusOK,
			response.GeneralResponse(
				"Student details retrieved successfully",
				http.StatusOK,
				student,
			),
		)

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

		response.WriteJson(w, http.StatusOK,
			response.GeneralResponse(
				"Students retrieved successfully",
				http.StatusOK,
				students,
			),
		)

	}
}

func DeleteStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("deleting student")
		id := r.PathValue("id")
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			// response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			response.WriteJson(w,
				http.StatusBadRequest,
				response.GeneralError(
					err,
					http.StatusBadRequest,
				),
			)

			return

		}

		slog.Info(fmt.Sprintf("id to be deleted %d", intId))
		deletedStudentId, err := storage.DeleteStudentById(intId)
		if err != nil {
			slog.Info("error while deleting student")

			if errors.Is(err, sql.ErrNoRows) {
				noStudentFoundErr := fmt.Errorf("no student found for the id %d", intId)
				// response.WriteJson(w, http.StatusNotFound, response.GeneralError(noStudentFoundErr))

				response.WriteJson(w,
					http.StatusNotFound,
					response.GeneralError(
						noStudentFoundErr,
						http.StatusNotFound,
					),
				)

				return

			}

			// response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			response.WriteJson(w,
				http.StatusInternalServerError,
				response.GeneralError(
					err,
					http.StatusInternalServerError,
				),
			)

			return
		}
		slog.Info(fmt.Sprintf("deleted student id %d", deletedStudentId))
		response.WriteJson(w, http.StatusOK,
			response.GeneralResponse(
				"Student deleted successfully",
				http.StatusOK,
				map[string]any{"message": "Student deleted successfully", "id": deletedStudentId},
			),
		)
	}
}

func UpdateStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("Updating student")

		var req model.StudentUpdateRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			if errors.Is(err, io.EOF) {

				response.WriteJson(w,
					http.StatusBadRequest,
					response.GeneralError(
						err,
						http.StatusBadRequest,
					),
				)

				return
			}

			response.WriteJson(w,
				http.StatusBadRequest,
				response.GeneralError(
					fmt.Errorf("error while parsing json"),
					http.StatusBadRequest,
				),
			)
			return
		}

		id := r.PathValue("id")
		studentId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w,
				http.StatusBadRequest,
				response.GeneralError(
					err,
					http.StatusBadRequest,
				),
			)

			return
		}

		updatedId, err := storage.UpdateStudentById(studentId, req)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				response.WriteJson(w,
					http.StatusNotFound,
					response.GeneralError(fmt.Errorf("no student found for this id"), http.StatusNotFound),
				)

				return

			}

			response.WriteJson(w,
				http.StatusInternalServerError,
				response.GeneralError(
					err,
					http.StatusInternalServerError,
				),
			)
			return
		}

		response.WriteJson(w, http.StatusOK,
			response.GeneralResponse(
				"Student updated successfully",
				http.StatusOK,
				map[string]any{
					"message": "success",
					"id":      updatedId,
				},
			),
		)

	}
}

func SearchStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("srarching student")

		query := r.URL.Query().Get("query")
		if strings.ToLower(query) == "" {

			emptyQueryErr := fmt.Errorf("please enter something to search")

			response.WriteJson(w,
				http.StatusBadRequest,
				response.GeneralError(
					emptyQueryErr,
					http.StatusBadRequest,
				),
			)
			return
		}

		students, err := storage.SearchStudent(query)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				response.WriteJson(w,
					http.StatusNotFound,
					response.GeneralError(fmt.Errorf("no student found for the given criteria"), http.StatusNotFound),
				)
				return
			}

			response.WriteJson(w,
				http.StatusInternalServerError,
				response.GeneralError(
					err,
					http.StatusInternalServerError,
				),
			)
			return

		}

		response.WriteJson(w, http.StatusOK,
			response.GeneralResponse(
				"Search completed successfully",
				http.StatusOK,
				students,
			),
		)

	}
}
