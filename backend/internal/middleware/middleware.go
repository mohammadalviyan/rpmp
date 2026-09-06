package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/mohammadalviyan/rpmp/backend/internal/domain"
)

const AccessCookieName = "rpmp_access"

type contextKey string

const (
	principalKey contextKey = "principal"
	requestIDKey contextKey = "request_id"
)

type TokenVerifier interface {
	Verify(string) (domain.Principal, error)
}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.NewString()
		w.Header().Set("X-Request-ID", requestID)
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestIDFromContext(ctx context.Context) string {
	requestID, _ := ctx.Value(requestIDKey).(string)
	return requestID
}

func RequireAuthentication(verifier TokenVerifier, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(AccessCookieName)
		if err != nil || cookie.Value == "" {
			writeUnauthenticated(w)
			return
		}
		principal, err := verifier.Verify(cookie.Value)
		if err != nil {
			writeUnauthenticated(w)
			return
		}
		ctx := context.WithValue(r.Context(), principalKey, principal)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ResolveAuthentication(verifier TokenVerifier, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(AccessCookieName)
		if err == nil && cookie.Value != "" {
			if principal, verifyErr := verifier.Verify(cookie.Value); verifyErr == nil {
				r = r.WithContext(context.WithValue(r.Context(), principalKey, principal))
			}
		}
		next.ServeHTTP(w, r)
	})
}

func PrincipalFromContext(ctx context.Context) (domain.Principal, bool) {
	principal, ok := ctx.Value(principalKey).(domain.Principal)
	return principal, ok
}

func RequireOrigin(allowedOrigin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if allowedOrigin == "" || r.Header.Get("Origin") != allowedOrigin {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"code": "forbidden", "message": "Request origin is not allowed.",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeUnauthenticated(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"code": "unauthenticated", "message": "Authentication is required.",
	})
}
