package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authDisabled := os.Getenv("AUTH_ENABLED") == "false" || os.Getenv("AUTH_ENABLED") == "0"
		if authDisabled {
			next.ServeHTTP(w, r)
			return
		}

		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		// Read the public key from the file
		publicKey, err := os.ReadFile(os.Getenv("SECNEX_GATEWAY_PUBLIC_KEY"))
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Parse the token
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		})
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !parsedToken.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
