package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github/com/ammar-nousher-ali/students-api/internal/model"
	"github/com/ammar-nousher-ali/students-api/internal/storage"
	"github/com/ammar-nousher-ali/students-api/internal/utils"
	"github/com/ammar-nousher-ali/students-api/internal/utils/response"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("mydevtestingkey123456789")

type SignUpRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"reqyured,min=6"`
	Role     string `json:"role" validate:"required,oneof=student teacher"`
}

func Signup(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req SignUpRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err, http.StatusBadRequest))
			return
		}

		req.Role = strings.ToLower(req.Role)
		exists, err := storage.IsEmailTaken(req.Email)
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

		userID, err := storage.CreateUser(user)
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
}

func SignIn(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var creds model.Creds

		err := json.NewDecoder(r.Body).Decode(&creds)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body"), http.StatusBadRequest))
			return
		}
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err, http.StatusBadRequest))
			return
		}

		error := validator.New().Struct(creds)
		if error != nil {
			var validation validator.ValidationErrors
			errors.As(error, &validation)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validation, http.StatusBadRequest))
			return

		}

		user, err := storage.GetUserByEmail(creds.Email)
		if err != nil {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(err, http.StatusUnauthorized))
			return
		}
		if !utils.CheckPasswordHash(creds.Password, user.Password) {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("invalid password"), http.StatusUnauthorized))
			return
		}

		expires := time.Now().Add(24 * time.Hour).Unix()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": user.ID,
			"email":   user.Email,
			"role":    user.Role,
			"expires": expires,
		})

		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err, http.StatusInternalServerError))
			return
		}

		response.WriteJson(
			w,
			http.StatusOK,
			response.GeneralResponse(
				"login successfull",
				http.StatusOK,
				map[string]any{
					"token":      tokenString,
					"expires_in": expires,
				},
			),
		)
	}
}
