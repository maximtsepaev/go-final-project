package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type signInInput struct {
	Password string `json:"password"`
}

func HandleSignIn(w http.ResponseWriter, r *http.Request) {
	var input signInInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		writeJSON(w, http.StatusOK, map[string]string{"token": "noauth"})
		return
	}

	// СДЕЛАТЬ НА РУССКОМ ОШИБКИ
	if input.Password != password {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid password"})
		return
	}

	passwordHash := sha256.Sum256([]byte(password))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": hex.EncodeToString(passwordHash[:]),
	})

	tokenString, err := token.SignedString([]byte(password))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to sign token"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}
