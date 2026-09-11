package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/mohammadalviyan/rpmp/backend/internal/domain"
	"github.com/mohammadalviyan/rpmp/backend/internal/middleware"
	"github.com/mohammadalviyan/rpmp/backend/internal/usecase"
)

type LoginUsecase interface {
	Execute(context.Context, usecase.LoginInput) (usecase.LoginOutput, error)
}

type CurrentUserUsecase interface {
	Execute(context.Context, domain.Principal) (domain.User, error)
}

type LogoutUsecase interface {
	Execute(context.Context, *domain.Principal, string) error
}

type CookieConfig struct {
	TTL    time.Duration
	Secure bool
}

type Auth struct {
	login   LoginUsecase
	current CurrentUserUsecase
	logout  LogoutUsecase
	cookie  CookieConfig
	now     func() time.Time
}

func NewAuth(login LoginUsecase, current CurrentUserUsecase, logout LogoutUsecase, cookie CookieConfig) *Auth {
	return &Auth{login: login, current: current, logout: logout, cookie: cookie, now: time.Now}
}

func (h *Auth) Routes(verifier middleware.TokenVerifier, allowedOrigin string) http.Handler {
	mux := http.NewServeMux()
	h.Mount(mux, verifier, allowedOrigin)
	return middleware.RequestID(middleware.RequestLog(nil, mux))
}

func (h *Auth) Mount(mux *http.ServeMux, verifier middleware.TokenVerifier, allowedOrigin string) {
	mux.Handle("POST /api/v1/auth/login", middleware.RequireOrigin(allowedOrigin, http.HandlerFunc(h.Login)))
	mux.Handle("GET /api/v1/auth/me", middleware.RequireAuthentication(verifier, http.HandlerFunc(h.Me)))
	logout := middleware.ResolveAuthentication(verifier, http.HandlerFunc(h.Logout))
	mux.Handle("POST /api/v1/auth/logout", middleware.RequireOrigin(allowedOrigin, logout))
}

func (h *Auth) Login(w http.ResponseWriter, r *http.Request) {
	var request struct {
		EmployeeID string `json:"employee_id"`
		Password   string `json:"password"`
	}
	if err := decodeJSON(w, r, &request); err != nil || request.EmployeeID == "" || request.Password == "" {
		writeError(w, domain.NewError(domain.KindInvalidRequest, err))
		return
	}
	output, err := h.login.Execute(r.Context(), usecase.LoginInput{
		EmployeeID: request.EmployeeID,
		Password:   request.Password,
		RequestID:  middleware.RequestIDFromContext(r.Context()),
	})
	if err != nil {
		writeError(w, err)
		return
	}
	http.SetCookie(w, h.accessCookie(output.Token, h.now().UTC().Add(h.cookie.TTL), int(h.cookie.TTL.Seconds())))
	writeJSON(w, http.StatusOK, map[string]any{"user": userResponse(output.User)})
}

func (h *Auth) Me(w http.ResponseWriter, r *http.Request) {
	principal, ok := middleware.PrincipalFromContext(r.Context())
	if !ok {
		writeError(w, domain.NewError(domain.KindUnauthenticated, nil))
		return
	}
	user, err := h.current.Execute(r.Context(), principal)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"user": userResponse(user)})
}

func (h *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	var principal *domain.Principal
	if resolved, ok := middleware.PrincipalFromContext(r.Context()); ok {
		principal = &resolved
	}
	err := h.logout.Execute(r.Context(), principal, middleware.RequestIDFromContext(r.Context()))
	http.SetCookie(w, h.accessCookie("", time.Unix(1, 0).UTC(), -1))
	if err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Auth) accessCookie(value string, expires time.Time, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name: middleware.AccessCookieName, Value: value, Path: "/",
		Expires: expires, MaxAge: maxAge, HttpOnly: true, Secure: h.cookie.Secure,
		SameSite: http.SameSiteLaxMode,
	}
}

type userJSON struct {
	ID          string      `json:"id"`
	EmployeeID  string      `json:"employee_id"`
	DisplayName string      `json:"display_name"`
	Role        domain.Role `json:"role"`
}

func userResponse(user domain.User) userJSON {
	return userJSON{ID: user.ID, EmployeeID: user.EmployeeID, DisplayName: user.DisplayName, Role: user.Role}
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain one JSON object")
	}
	return nil
}

func writeError(w http.ResponseWriter, err error) {
	kind := domain.ErrorKindOf(err)
	status := http.StatusInternalServerError
	message := "An internal error occurred."
	switch kind {
	case domain.KindInvalidRequest:
		status, message = http.StatusBadRequest, "Request body is invalid."
	case domain.KindInvalidCredentials:
		status, message = http.StatusUnauthorized, "Employee ID or password is incorrect."
	case domain.KindAccountInactive:
		status, message = http.StatusForbidden, "Account is inactive."
	case domain.KindUnauthenticated:
		status, message = http.StatusUnauthorized, "Authentication is required."
	case domain.KindForbidden:
		status, message = http.StatusForbidden, "You are not authorized to perform this action."
	case domain.KindInvalidPeriod:
		status, message = http.StatusBadRequest, "The requested period is invalid."
	case domain.KindSourceUnavailable:
		status, message = http.StatusServiceUnavailable, "Operational data is temporarily unavailable."
	case domain.KindSyncInProgress:
		status, message = http.StatusConflict, "A synchronization is already in progress."
	default:
		kind = domain.KindInternal
	}
	writeJSON(w, status, map[string]string{"code": string(kind), "message": message})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
