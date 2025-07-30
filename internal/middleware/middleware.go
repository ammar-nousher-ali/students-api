package middleware

import (
	"fmt"
	"github/com/ammar-nousher-ali/students-api/internal/http/handlers/auth"
	"github/com/ammar-nousher-ali/students-api/internal/utils/response"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			err := fmt.Errorf("missing authorization header")
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(err, http.StatusUnauthorized))
			return
		}
		//slog.Info("token is ", authHeader)
		parts := strings.Split(authHeader, " ")
		//slog.Info("parts is ", parts)
		if len(parts) != 2 || parts[0] != "Bearer" {
			err := fmt.Errorf("invalid authorization header format")
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(err, http.StatusUnauthorized))
			return
		}

		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC) //This is type assertion. it means I expect token.Method to be of type *jwt.SigningMethodHMAC. Try to cast it to that type.
			if !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return auth.JwtSecret, nil
		})

		if err != nil || !token.Valid {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("invalid or expired token"), http.StatusUnauthorized))
			return
		}

		next(w, r)
	}
}
