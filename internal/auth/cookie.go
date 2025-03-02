package auth

import (
	"context"
	"net/http"
	"time"
)

type contextKey string

const userIDKey = contextKey("userID")

func AutoAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := GetTokenFromCookie(r)
		userID, err := GetUserID(tokenString)
		if err != nil || userID == "" {
			newToken, err := BuildJWTString()
			if err != nil {
				http.Error(w, "Failed to generate token", http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:     "jwt",
				Value:    newToken,
				Expires:  time.Now().Add(time.Hour),
				HttpOnly: true,
				Path:     "/",
			})
			userID, err = GetUserID(newToken)
			if err != nil {
				http.Error(w, "Failed to get user ID", http.StatusInternalServerError)
				return
			}
		}
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}

func StrictAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := GetTokenFromCookie(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userID, err := GetUserID(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}

func GetTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func GetUserIDFromContext(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(userIDKey).(string)
	return userID, ok
}
