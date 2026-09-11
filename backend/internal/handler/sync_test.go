package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mohammadalviyan/rpmp/backend/internal/domain"
	"github.com/mohammadalviyan/rpmp/backend/internal/usecase"
)

type syncStatusRepositoryStub struct {
	run domain.SyncRun
	err error
}

func (s syncStatusRepositoryStub) LatestSyncRun(context.Context) (domain.SyncRun, error) {
	return s.run, s.err
}

type backgroundSyncRunnerStub struct {
	run domain.SyncRun
	err error
}

func (s *backgroundSyncRunnerStub) Start(context.Context, context.Context) (domain.SyncRun, error) {
	return s.run, s.err
}

func syncRouter(role domain.Role, status usecase.SyncStatusRepository, runner usecase.BackgroundSyncRunner) http.Handler {
	syncHandler := NewSync(
		usecase.NewGetSyncStatus(status),
		usecase.NewStartSync(runner, context.Background()),
	)
	return NewRouter(
		NewAuth(nil, nil, nil, CookieConfig{}),
		nil,
		verifierStub{principal: domain.Principal{UserID: "user-1", Role: role}},
		"https://rpmp.example",
		syncHandler,
	)
}

func authenticatedRequest(method, path string) *http.Request {
	request := httptest.NewRequest(method, path, nil)
	request.AddCookie(&http.Cookie{Name: "rpmp_access", Value: "valid"})
	return request
}

func TestSyncStatusContractForViewerAndAdmin(t *testing.T) {
	startedAt := time.Date(2026, 9, 9, 0, 0, 0, 0, time.UTC)
	finishedAt := startedAt.Add(time.Minute)
	sourceDetail := "source CSV /private/path failed"
	repository := syncStatusRepositoryStub{run: domain.SyncRun{
		ID: "018f5f71-9cb9-7a61-97e7-d8f33f12b821", Status: domain.SyncRunFailure,
		StartedAt: startedAt, FinishedAt: &finishedAt, RowsWritten: 100,
		ErrorCode: &sourceDetail,
	}}

	for _, role := range []domain.Role{domain.RoleViewer, domain.RoleAdmin} {
		request := authenticatedRequest(http.MethodGet, "/api/v1/sync/status")
		response := httptest.NewRecorder()
		syncRouter(role, repository, &backgroundSyncRunnerStub{}).ServeHTTP(response, request)
		if response.Code != http.StatusOK {
			t.Fatalf("role %s status = %d %s", role, response.Code, response.Body.String())
		}
		body := response.Body.String()
		for _, want := range []string{
			`"run_id":"018f5f71-9cb9-7a61-97e7-d8f33f12b821"`,
			`"status":"failure"`,
			`"started_at":"2026-09-09T00:00:00Z"`,
			`"finished_at":"2026-09-09T00:01:00Z"`,
			`"rows_written":100`,
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("role %s missing %s in %s", role, want, body)
			}
		}
		if strings.Contains(body, sourceDetail) || strings.Contains(body, "/private/path") {
			t.Fatalf("source detail leaked: %s", body)
		}
	}
}

func TestSyncStatusReturnsNeverWithNullFields(t *testing.T) {
	router := syncRouter(
		domain.RoleViewer,
		syncStatusRepositoryStub{err: domain.NewError(domain.KindNotFound, errors.New("no rows"))},
		&backgroundSyncRunnerStub{},
	)
	request := authenticatedRequest(http.MethodGet, "/api/v1/sync/status")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	var body map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if response.Code != http.StatusOK || body["status"] != "never" ||
		body["run_id"] != nil || body["started_at"] != nil || body["finished_at"] != nil ||
		body["rows_written"] != float64(0) {
		t.Fatalf("status=%d body=%#v", response.Code, body)
	}
}

func TestAdminStartsSyncAndOverlapReturnsConflict(t *testing.T) {
	runner := &backgroundSyncRunnerStub{run: domain.SyncRun{
		ID: "018f5f71-9cb9-7a61-97e7-d8f33f12b822", Status: domain.SyncRunRunning,
	}}
	router := syncRouter(domain.RoleAdmin, syncStatusRepositoryStub{}, runner)
	request := authenticatedRequest(http.MethodPost, "/api/v1/sync")
	request.Header.Set("Origin", "https://rpmp.example")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusAccepted ||
		response.Body.String() != "{\"run_id\":\"018f5f71-9cb9-7a61-97e7-d8f33f12b822\",\"status\":\"running\"}\n" {
		t.Fatalf("accepted = %d %s", response.Code, response.Body.String())
	}

	runner.err = domain.NewError(domain.KindSyncInProgress, errors.New("advisory lock held"))
	request = authenticatedRequest(http.MethodPost, "/api/v1/sync")
	request.Header.Set("Origin", "https://rpmp.example")
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusConflict ||
		!strings.Contains(response.Body.String(), `"code":"sync_in_progress"`) ||
		strings.Contains(response.Body.String(), "advisory") {
		t.Fatalf("conflict = %d %s", response.Code, response.Body.String())
	}
}

func TestViewerCannotStartSync(t *testing.T) {
	runner := &backgroundSyncRunnerStub{}
	request := authenticatedRequest(http.MethodPost, "/api/v1/sync")
	request.Header.Set("Origin", "https://rpmp.example")
	response := httptest.NewRecorder()
	syncRouter(domain.RoleViewer, syncStatusRepositoryStub{}, runner).ServeHTTP(response, request)
	if response.Code != http.StatusForbidden ||
		!strings.Contains(response.Body.String(), `"code":"forbidden"`) {
		t.Fatalf("viewer = %d %s", response.Code, response.Body.String())
	}
}
