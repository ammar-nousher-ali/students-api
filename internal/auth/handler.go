package auth

import (
	"encoding/json"
	"fmt"
	"github/com/ammar-nousher-ali/students-api/internal/model"
	"github/com/ammar-nousher-ali/students-api/internal/storage"
	"github/com/ammar-nousher-ali/students-api/internal/utils/response"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type SignUpRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"reqyured,min=6"`
	Role     string `json:"role" validate:"required,oneof=student teach admin"`
}

type AuthHandler struct {
	Storage storage.Storage
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err, http.StatusBadRequest))
		return
	}

	req.Role = strings.ToLower(req.Role)
	exists, err := h.Storage.IsEmailTaken(req.Email)
	if err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err, http.StatusInternalServerError))
		return
	}

	if exists {
		response.WriteJson(w, http.StatusConflict,
			response.GeneralError(
				fmt.Errorf("user with email %s already exists. please try again with different email", req.Email),
				http.StatusConflict),
		)

		return
	}

	//hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(
			fmt.Errorf("something went wrong"),
			http.StatusInternalServerError,
		),
		)
	}

	user := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
	}

	userID, err := h.Storage.CreateUser(user)
	if err != nil {
		response.WriteJson(w, http.StatusInternalServerError,
			response.GeneralResponse(
				"error while creating user",
				http.StatusInternalServerError,
				response.GeneralError(err,
					http.StatusInternalServerError),
			),
		)
	}

	user.ID = userID
	user.Password = ""

	response.WriteJson(w, http.StatusCreated,
		response.GeneralResponse(
			"user created successfully",
			http.StatusCreated,
			map[string]any{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
				"role":  user.Role,
			},
		),
	)
}
