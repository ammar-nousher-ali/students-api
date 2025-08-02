package course

import (
	"encoding/json"
	"errors"
	"fmt"
	"github/com/ammar-nousher-ali/students-api/internal/model"
	"github/com/ammar-nousher-ali/students-api/internal/storage"
	"github/com/ammar-nousher-ali/students-api/internal/utils/response"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var course model.Course

		slog.Info("creating student")

		err := json.NewDecoder(r.Body).Decode(&course)
		if err != nil {
			if errors.Is(err, io.EOF) {
				response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body"), http.StatusBadRequest))
				return
			}

		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request"), http.StatusBadRequest))
			return
		}

		now := time.Now()
		course.UpdatedAt = now
		course.CreatedAt = now

		error := validator.New().Struct(course)

		if error != nil {
			var validation validator.ValidationErrors
			errors.As(error, validation)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validation, http.StatusBadRequest))
			return
		}

		id, err := storage.CreateCourse(course)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err, http.StatusBadRequest))
			return
		}

		course.Id = id
		response.WriteJson(w, http.StatusOK, response.GeneralResponse("course create successfully", http.StatusOK, course))

	}
}

func NewBatch(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var courses []model.Course

		err := json.NewDecoder(r.Body).Decode(&courses)

		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body"), http.StatusBadRequest))
			return

		}
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err, http.StatusBadRequest))
			return
		}

		var batchResponse response.BatchResponse
		for _, course := range courses {
			now := time.Now()
			course.UpdatedAt = now
			course.CreatedAt = now
			id, err := storage.CreateCourse(course)
			if err != nil {

				var reason string
				if strings.Contains(err.Error(), "UNIQUE constraint failed") {
					reason = "course already added"
				} else {
					reason = err.Error()
				}
				batchResponse.Data = append(batchResponse.Data, response.BatchData{
					Success: false,
					Data: map[string]any{
						"message": "failed",
						"reason":  reason,
					},
				})
			} else {
				course.Id = id
				batchResponse.Data = append(batchResponse.Data, response.BatchData{
					Success: true,
					Data: map[string]any{
						"course": course,
					},
				})

			}

		}

		response.WriteJson(w, http.StatusOK, response.GeneralBatchResponse("success", http.StatusOK, batchResponse.Data))

	}

}
