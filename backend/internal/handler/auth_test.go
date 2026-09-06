package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mohammadalviyan/rpmp/backend/internal/domain"
	"github.com/mohammadalviyan/rpmp/backend/internal/usecase"
)

type loginStub struct {
	output usecase.LoginOutput
	err    error
}

func (s loginStub) Execute(context.Context, usecase.LoginInput) (usecase.LoginOutput, error) {
	return s.output, s.err
}

type currentStub struct {
	user domain.User
	err  error
}

func (s currentStub) Execute(context.Context, domain.Principal) (domain.User, error) {
	return s.user, s.err
}

type logoutStub struct {
	principal *domain.Principal
	called    bool
}

func (s *logoutStub) Execute(_ context.Context, principal *domain.Principal, _ string) error {
	s.called = true
	s.principal = principal
	return nil
}

type verifierStub struct {
	principal domain.Principal
	err       error
}

func (s verifierStub) Verify(string) (domain.Principal, error) {
	return s.principal, s.err
}

func frozenUser() domain.User {
	return domain.User{
		ID: "018f5f71-9cb9-7a61-97e7-d8f33f12b821", EmployeeID: "12345678",
		DisplayName: "Example Viewer", PasswordHash: "must-not-leak", Role: domain.RoleViewer, Active: true,
	}
}

func TestLoginContractAndCookie(t *testing.T) {
	user := frozenUser()
	logout := &logoutStub{}
	h := NewAuth(
		loginStub{output: usecase.LoginOutput{User: user, Token: "jwt-secret-value"}},
		currentStub{user: user},
		logout,
		CookieConfig{TTL: 30 * time.Minute, Secure: true},
	)
	h.now = func() time.Time { return time.Date(2026, 9, 6, 7, 0, 0, 0, time.UTC) }
	router := h.Routes(verifierStub{}, "https://rpmp.example")

	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(`{"employee_id":"12345678","password":"secret"}`))
	request.Header.Set("Origin", "https://rpmp.example")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if strings.Contains(response.Body.String(), "secret") || strings.Contains(response.Body.String(), "password") {
		t.Fatalf("secret material leaked: %s", response.Body.String())
	}
	if response.Body.String() != `{"user":{"id":"018f5f71-9cb9-7a61-97e7-d8f33f12b821","employee_id":"12345678","display_name":"Example Viewer","role":"viewer"}}`+"\n" {
		t.Fatalf("body = %q", response.Body.String())
	}
	cookies := response.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("cookies = %#v", cookies)
	}
	cookie := cookies[0]
	if cookie.Name != "rpmp_access" || cookie.Value != "jwt-secret-value" || !cookie.HttpOnly || !cookie.Secure ||
		cookie.Path != "/" || cookie.Domain != "" || cookie.SameSite != http.SameSiteLaxMode || cookie.MaxAge != 1800 {
		t.Fatalf("unexpected cookie: %#v", cookie)
	}
}

func TestLoginRejectsMalformedBodyAndOrigin(t *testing.T) {
	h := NewAuth(loginStub{}, currentStub{}, &logoutStub{}, CookieConfig{TTL: 30 * time.Minute})
	router := h.Routes(verifierStub{}, "http://localhost:3000")

	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(`{"employee_id":"1","password":`))
	request.Header.Set("Origin", "http://localhost:3000")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest || !strings.Contains(response.Body.String(), `"code":"invalid_request"`) {
		t.Fatalf("malformed response = %d %s", response.Code, response.Body.String())
	}

	request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(`{"employee_id":"1","password":"x"}`))
	request.Header.Set("Origin", "https://evil.example")
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusForbidden {
		t.Fatalf("origin status = %d", response.Code)
	}
}

func TestLoginMapsSafeCredentialError(t *testing.T) {
	h := NewAuth(
		loginStub{err: domain.NewError(domain.KindInvalidCredentials, errors.New("database detail"))},
		currentStub{}, &logoutStub{}, CookieConfig{TTL: 30 * time.Minute},
	)
	router := h.Routes(verifierStub{}, "https://rpmp.example")
	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(`{"employee_id":"1","password":"top-secret"}`))
	request.Header.Set("Origin", "https://rpmp.example")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusUnauthorized ||
		response.Body.String() != "{\"code\":\"invalid_credentials\",\"message\":\"Employee ID or password is incorrect.\"}\n" {
		t.Fatalf("response = %d %q", response.Code, response.Body.String())
	}
}

func TestMeAndLogoutContracts(t *testing.T) {
	user := frozenUser()
	logout := &logoutStub{}
	principal := domain.Principal{UserID: user.ID, Role: domain.RoleViewer}
	h := NewAuth(loginStub{}, currentStub{user: user}, logout, CookieConfig{TTL: 30 * time.Minute, Secure: true})
	router := h.Routes(verifierStub{principal: principal}, "https://rpmp.example")

	request := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	request.AddCookie(&http.Cookie{Name: "rpmp_access", Value: "valid"})
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"employee_id":"12345678"`) {
		t.Fatalf("me response = %d %s", response.Code, response.Body.String())
	}

	request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	request.Header.Set("Origin", "https://rpmp.example")
	request.AddCookie(&http.Cookie{Name: "rpmp_access", Value: "valid"})
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent || response.Body.Len() != 0 || !logout.called || logout.principal == nil {
		t.Fatalf("logout response = %d %q, stub=%#v", response.Code, response.Body.String(), logout)
	}
	cookie := response.Result().Cookies()[0]
	if cookie.Name != "rpmp_access" || cookie.MaxAge != -1 || !cookie.HttpOnly || !cookie.Secure ||
		cookie.Path != "/" || cookie.Domain != "" || cookie.SameSite != http.SameSiteLaxMode {
		t.Fatalf("logout cookie = %#v", cookie)
	}
}
