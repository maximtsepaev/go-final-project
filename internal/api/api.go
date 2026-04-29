// Package api содержит HTTP-обработчики и маршруты планировщика задач.
package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

const webDir = "./web" // Путь к директории с веб-ресурсами (HTML, CSS, JS)

// DateFormat задает формат даты для хранения и обработки (ГГГГММДД).
const DateFormat = "20060102"

// Init регистрирует маршруты API и раздачу статических файлов.
func Init(r chi.Router) {
	r.Get("/api/nextdate", HandleNextDate)
	r.With(auth).Get("/api/tasks", HandleGetTasks)
	r.Post("/api/signin", HandleSignIn)

	// Роут для разных методов на одном пути
	r.Route("/api/task", func(r chi.Router) {
		r.Use(auth)
		r.Post("/", HandleAddTask)
		r.Get("/", HandleGetTaskByID)
		r.Put("/", HandleUpdateTask)
		r.Post("/done", HandleMarkTaskAsDone)
		r.Delete("/", HandleDeleteTask)
	})

	// Раздача статических файлов из директории web
	r.Handle("/*", http.FileServer(http.Dir(webDir)))
}

// writeJSON отправляет JSON-ответ с указанным статусом и данными.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// auth является middleware для проверки JWT-токена в cookie и авторизации пользователя.
func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")

		if len(pass) > 0 {
			var token string
			var valid bool

			cookie, err := r.Cookie("token")
			if err == nil {
				token = cookie.Value
			}

			expectedHash := sha256.Sum256([]byte(pass))
			expectedHashHex := hex.EncodeToString(expectedHash[:])

			claims := jwt.MapClaims{}
			jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("Invalid signing method")
				}
				return []byte(pass), nil
			})
			if err == nil && jwtToken.Valid {
				hashClaim, ok := claims["hash"].(string)
				if ok && hashClaim == expectedHashHex {
					valid = true
				}
			}

			if !valid {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
