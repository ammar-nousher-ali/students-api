package enroll_student

import (
	"encoding/json"
	"fmt"
	"github/com/ammar-nousher-ali/students-api/internal/model"
	"github/com/ammar-nousher-ali/students-api/internal/storage"
	"github/com/ammar-nousher-ali/students-api/internal/utils/response"
	"net/http"
	"strconv"
)

func EnrollStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		strId := r.PathValue("student_id")
		if len(strId) == 0 {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("please add correct path params"), http.StatusBadRequest))
			return

		}

		studentId, err := strconv.ParseInt(strId, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid student id"), http.StatusBadRequest))
			return

		}

		var req model.EnrollRequest
		error := json.NewDecoder(r.Body).Decode(&req)
		if error != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request"), http.StatusBadRequest))
			return
		}

		result, err := storage.EnrollStudentInCourse(studentId, req)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err, http.StatusInternalServerError))
			return
		}

		response.WriteJson(w, http.StatusOK, response.GeneralResponse("success", http.StatusOK, result))

	}
}

func GetStudentWithEnrolledCourse(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		strId := r.PathValue("student_id")
		if len(strId) == 0 {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("please add correct path params"), http.StatusBadRequest))
			return
		}

		studentId, err := strconv.ParseInt(strId, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid student id"), http.StatusBadRequest))
			return
		}

		studentWithCoursesResponse, err := storage.FetchStudentWithEnrolledCourse(studentId)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err, http.StatusInternalServerError))
			return
		}

		response.WriteJson(w, http.StatusOK, response.GeneralResponse("success", http.StatusOK, studentWithCoursesResponse))

	}
}
