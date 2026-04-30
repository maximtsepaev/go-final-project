// Package api содержит HTTP-обработчики и маршруты планировщика задач.
package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

const webDir = "./web" // Путь к директории с веб-ресурсами (HTML, CSS, JS)

// DateFormat задает формат даты для хранения и обработки (ГГГГММДД).
const DateFormat = "20060102"

var password string     // Значение TODO_PASSWORD
var passwordHash string // SHA256 хэш пароля

// Init регистрирует маршруты API и раздачу статических файлов.
func Init(r chi.Router) {
	slog.Info("routes initialized", "static_dir", webDir)

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

// init выполняется при загрузке пакета и вычисляет хэш пароля для генерации JWT-токена,
// а также для его проверки в middleware auth, но только в том случае,
// если переменная окружения TODO_PASSWORD установлена.
func init() {
	password = os.Getenv("TODO_PASSWORD")
	if password != "" {
		hash := sha256.Sum256([]byte(password))
		passwordHash = hex.EncodeToString(hash[:])
	}
}

// auth является middleware для проверки JWT-токена в cookie и авторизации пользователя.
func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string
		var valid bool
		if len(password) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value == "" {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "authentication required"})
			return
		}

		token = cookie.Value

		claims := jwt.MapClaims{}
		jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return []byte(password), nil
		})
		if err == nil && jwtToken.Valid {
			exp, ok := claims["exp"].(float64)
			if !ok || time.Now().Unix() > int64(exp) {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "token expired"})
				return
			}

			hashClaim, ok := claims["hash"].(string)
			if ok && hashClaim == passwordHash {
				valid = true
			}
		}

		if !valid {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "authentication required"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
