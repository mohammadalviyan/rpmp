package handler

import (
	"net/http"

	"github.com/mohammadalviyan/rpmp/backend/internal/middleware"
)

func NewRouter(
	auth *Auth,
	dashboard *Dashboard,
	verifier middleware.TokenVerifier,
	allowedOrigin string,
	syncHandler *Sync,
	useCases ...*UseCases,
) http.Handler {
	mux := http.NewServeMux()
	auth.Mount(mux, verifier, allowedOrigin)
	if dashboard != nil {
		dashboard.Mount(mux, verifier)
	}
	if syncHandler != nil {
		syncHandler.Mount(mux, verifier, allowedOrigin)
	}
	if len(useCases) > 0 && useCases[0] != nil {
		useCases[0].Mount(mux, verifier)
	}
	return middleware.RequestID(middleware.RequestLog(nil, mux))
}
