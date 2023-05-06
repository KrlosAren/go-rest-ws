package middleware

import (
	"net/http"
	"strings"

	"github.com/go/rest-ws/models"
	"github.com/go/rest-ws/server"
	"github.com/golang-jwt/jwt"
)

var (
	NO_AUTH_NEEDED = []string{
		"login",
		"signup",
	}
)

func shouldCheckToken(route string) bool {
	for _, p := range NO_AUTH_NEEDED {
		if strings.Contains(route, p) {
			return false
		}
	}
	return true
}

func CheckAuthMiddleware(s server.Server) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !shouldCheckToken(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}
			tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
			token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(s.Config().JWTSecret), nil
			})

			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				r.Header.Set("userId", claims.UserId)
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}
		})

	}
}
