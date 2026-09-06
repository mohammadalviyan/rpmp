package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mohammadalviyan/rpmp/backend/internal/domain"
)

type fakeVerifier struct {
	principal domain.Principal
	err       error
}

func (f fakeVerifier) Verify(string) (domain.Principal, error) {
	return f.principal, f.err
}

func TestRequireAuthentication(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		principal, ok := PrincipalFromContext(r.Context())
		if !ok || principal.UserID != "user-1" {
			t.Fatalf("principal missing: %#v", principal)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	handler := RequireAuthentication(fakeVerifier{
		principal: domain.Principal{UserID: "user-1", Role: domain.RoleViewer},
	}, next)

	request := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	request.AddCookie(&http.Cookie{Name: AccessCookieName, Value: "token"})
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d", response.Code)
	}
}

func TestRequireAuthenticationRejectsMissingAndInvalidCookie(t *testing.T) {
	for _, test := range []struct {
		name   string
		cookie *http.Cookie
	}{
		{name: "missing"},
		{name: "invalid", cookie: &http.Cookie{Name: AccessCookieName, Value: "bad"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			handler := RequireAuthentication(fakeVerifier{err: errors.New("invalid")}, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
				t.Fatal("protected handler called")
			}))
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			if test.cookie != nil {
				request.AddCookie(test.cookie)
			}
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != http.StatusUnauthorized || response.Body.String() != "{\"code\":\"unauthenticated\",\"message\":\"Authentication is required.\"}\n" {
				t.Fatalf("response = %d %q", response.Code, response.Body.String())
			}
		})
	}
}

func TestRequireOrigin(t *testing.T) {
	handler := RequireOrigin("https://rpmp.example", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	for _, test := range []struct {
		origin string
		want   int
	}{
		{origin: "https://rpmp.example", want: http.StatusNoContent},
		{origin: "https://evil.example", want: http.StatusForbidden},
		{origin: "", want: http.StatusForbidden},
	} {
		request := httptest.NewRequest(http.MethodPost, "/", nil)
		request.Header.Set("Origin", test.origin)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != test.want {
			t.Fatalf("origin %q status = %d, want %d", test.origin, response.Code, test.want)
		}
	}
}

func TestResolveAuthenticationIgnoresInvalidToken(t *testing.T) {
	handler := ResolveAuthentication(fakeVerifier{err: errors.New("invalid")}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := PrincipalFromContext(r.Context()); ok {
			t.Fatal("invalid identity resolved")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	request := httptest.NewRequest(http.MethodPost, "/", nil)
	request.AddCookie(&http.Cookie{Name: AccessCookieName, Value: "bad"})
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d", response.Code)
	}
}
