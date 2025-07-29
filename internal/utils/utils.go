package utils

import "golang.org/x/crypto/bcrypt"

func CheckPasswordHash(password string, hash string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil

}
