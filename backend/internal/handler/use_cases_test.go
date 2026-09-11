package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mohammadalviyan/rpmp/backend/internal/domain"
	"github.com/mohammadalviyan/rpmp/backend/internal/usecase"
)

type listUseCasesStub struct {
	result domain.UseCaseList
	err    error
	input  usecase.ListUseCasesInput
}

func (s *listUseCasesStub) Execute(
	_ context.Context,
	input usecase.ListUseCasesInput,
) (domain.UseCaseList, error) {
	s.input = input
	return s.result, s.err
}

type getUseCaseStub struct {
	result domain.UseCaseDetail
	err    error
	input  usecase.GetUseCaseInput
}

func (s *getUseCaseStub) Execute(
	_ context.Context,
	input usecase.GetUseCaseInput,
) (domain.UseCaseDetail, error) {
	s.input = input
	return s.result, s.err
}

func TestUseCaseListFrozenResponseAndCookieAuth(t *testing.T) {
	freshAt := time.Date(2026, 9, 11, 2, 0, 0, 0, time.UTC)
	rate := 96.15
	list := &listUseCasesStub{result: domain.UseCaseList{
		Freshness: domain.Freshness{Status: domain.FreshnessFresh, LastSuccessfulRefreshAt: freshAt},
		Items: []domain.UseCaseMetrics{
			{
				ID: "018f5f71-9cb9-7a61-97e7-d8f33f12b821", Name: "Finance/Invoice.Process",
				Status: domain.UseCaseActive, ProcessCount: 2, SuccessfulCount: 100,
				ErrorCount: 3, StoppedCount: 1, FailedCount: 4, ExecutionVolume: 104,
				SuccessRate: &rate, Environments: []string{"PROD"},
			},
		},
	}}
	router := NewRouter(
		NewAuth(loginStub{}, currentStub{}, &logoutStub{}, CookieConfig{}),
		nil,
		verifierStub{principal: domain.Principal{UserID: "viewer-1", Role: domain.RoleViewer}},
		"https://rpmp.example",
		nil,
		NewUseCases(list, &getUseCaseStub{}),
	)

	request := httptest.NewRequest(http.MethodGet, "/api/v1/use-cases?q=invoice&status=active", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated = %d %s", response.Code, response.Body.String())
	}

	request = httptest.NewRequest(http.MethodGet, "/api/v1/use-cases?q=invoice&status=active", nil)
	request.AddCookie(&http.Cookie{Name: "rpmp_access", Value: "valid"})
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d %s", response.Code, response.Body.String())
	}
	if list.input.Query != "invoice" || list.input.Status != "active" ||
		list.input.Principal.Role != domain.RoleViewer {
		t.Fatalf("input = %#v", list.input)
	}
	var body map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	freshness := body["freshness"].(map[string]any)
	items := body["items"].([]any)
	item := items[0].(map[string]any)
	if freshness["status"] != "fresh" ||
		freshness["last_successful_refresh_at"] != "2026-09-11T02:00:00Z" ||
		item["failed_count"] != float64(4) || item["execution_volume"] != float64(104) ||
		item["success_rate"] != 96.15 || item["source_key"] != nil || item["processes"] != nil {
		t.Fatalf("body = %#v", body)
	}
}

func TestUseCaseDetailFrozenResponse(t *testing.T) {
	freshAt := time.Date(2026, 9, 11, 2, 0, 0, 0, time.UTC)
	environment := "PROD"
	rate := 66.67
	get := &getUseCaseStub{result: domain.UseCaseDetail{
		Freshness: domain.Freshness{Status: domain.FreshnessFresh, LastSuccessfulRefreshAt: freshAt},
		UseCase: domain.UseCaseMetrics{
			ID: "018f5f71-9cb9-7a61-97e7-d8f33f12b821", Name: "Finance/Invoice.Process",
			Status: domain.UseCaseActive, SourceKey: "Finance/Invoice.Process", UpdatedAt: freshAt,
			ProcessCount: 1, SuccessfulCount: 2, ErrorCount: 1, FailedCount: 1,
			ExecutionVolume: 3, SuccessRate: &rate, Environments: []string{"PROD"},
			Processes: []domain.UseCaseProcessSnapshot{
				{
					SourceProcessKey: "42", ProcessName: "Invoice", PackageName: "Invoice.1.0.0",
					EnvironmentName: &environment, SuccessfulCount: 2, ErrorCount: 1,
				},
			},
		},
	}}
	router := NewRouter(
		NewAuth(loginStub{}, currentStub{}, &logoutStub{}, CookieConfig{}),
		nil,
		verifierStub{principal: domain.Principal{UserID: "admin-1", Role: domain.RoleAdmin}},
		"https://rpmp.example",
		nil,
		NewUseCases(&listUseCasesStub{}, get),
	)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/use-cases/018f5f71-9cb9-7a61-97e7-d8f33f12b821", nil)
	request.AddCookie(&http.Cookie{Name: "rpmp_access", Value: "valid"})
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d %s", response.Code, response.Body.String())
	}
	body := response.Body.String()
	for _, field := range []string{
		`"source_key":"Finance/Invoice.Process"`,
		`"updated_at":"2026-09-11T02:00:00Z"`,
		`"source_process_key":"42"`,
		`"failed_count":1`,
		`"executing_count":0`,
	} {
		if !strings.Contains(body, field) {
			t.Fatalf("missing %s in %s", field, body)
		}
	}
	if get.input.Principal.Role != domain.RoleAdmin {
		t.Fatalf("input = %#v", get.input)
	}
}

func TestUseCaseErrorsAndNeverSyncedShape(t *testing.T) {
	for _, test := range []struct {
		name string
		path string
		list *listUseCasesStub
		get  *getUseCaseStub
		code int
		want string
	}{
		{
			name: "invalid status", path: "/api/v1/use-cases?status=paused",
			list: &listUseCasesStub{err: domain.NewError(domain.KindInvalidRequest, nil)},
			get:  &getUseCaseStub{}, code: http.StatusBadRequest, want: `"code":"invalid_request"`,
		},
		{
			name: "malformed id", path: "/api/v1/use-cases/not-a-uuid",
			list: &listUseCasesStub{},
			get:  &getUseCaseStub{err: domain.NewError(domain.KindNotFound, nil)},
			code: http.StatusNotFound, want: `"code":"not_found"`,
		},
		{
			name: "unknown id", path: "/api/v1/use-cases/018f5f71-9cb9-7a61-97e7-d8f33f12b822",
			list: &listUseCasesStub{},
			get:  &getUseCaseStub{err: domain.NewError(domain.KindNotFound, nil)},
			code: http.StatusNotFound, want: `"code":"not_found"`,
		},
		{
			name: "never synced", path: "/api/v1/use-cases",
			list: &listUseCasesStub{result: domain.UseCaseList{
				Freshness: domain.Freshness{Status: domain.FreshnessNever},
				Items:     []domain.UseCaseMetrics{},
			}},
			get: &getUseCaseStub{}, code: http.StatusOK, want: `"freshness":{"status":"never"}`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			router := NewRouter(
				NewAuth(loginStub{}, currentStub{}, &logoutStub{}, CookieConfig{}),
				nil,
				verifierStub{principal: domain.Principal{Role: domain.RoleViewer}},
				"https://rpmp.example",
				nil,
				NewUseCases(test.list, test.get),
			)
			request := httptest.NewRequest(http.MethodGet, test.path, nil)
			request.AddCookie(&http.Cookie{Name: "rpmp_access", Value: "valid"})
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			if response.Code != test.code || !strings.Contains(response.Body.String(), test.want) {
				t.Fatalf("response = %d %s", response.Code, response.Body.String())
			}
			if test.name == "never synced" &&
				strings.Contains(response.Body.String(), "last_successful_refresh_at") {
				t.Fatalf("unexpected timestamp: %s", response.Body.String())
			}
		})
	}
}
