package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// signInInput представляет структуру данных для входа пользователя.
type signInInput struct {
	Password string `json:"password"`
}

// HandleSignIn проверяет пароль и возвращает JWT-токен.
func HandleSignIn(w http.ResponseWriter, r *http.Request) {
	var input signInInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	password := os.Getenv("TODO_PASSWORD")
	if len(password) == 0 {
		writeJSON(w, http.StatusOK, map[string]string{"token": "noauth"})
		return
	}

	if input.Password != password {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": passwordHash,
		"exp":  time.Now().Add(8 * time.Hour).Unix(), // Для безопасности, несмотря на то, что фронт устанавливает лайф тайм куки
	})

	tokenString, err := token.SignedString([]byte(password))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to sign token"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}
