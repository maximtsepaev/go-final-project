package api

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

const webDir = "./web"        // Путь к директории с веб-ресурсами (HTML, CSS, JS)
const DateFormat = "20060102" // Формат даты для хранения и обработки (ГГГГММДД)

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

	r.Handle("/*", http.FileServer(http.Dir(webDir)))
}

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
				// возвращаем ошибку авторизации 401
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
