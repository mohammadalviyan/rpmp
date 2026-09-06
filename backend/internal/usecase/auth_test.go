package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/mohammadalviyan/rpmp/backend/internal/domain"
)

type fakeUsers struct {
	byEmployee domain.User
	byID       domain.User
	err        error
}

func (f *fakeUsers) FindByEmployeeID(context.Context, string) (domain.User, error) {
	return f.byEmployee, f.err
}

func (f *fakeUsers) FindByID(context.Context, string) (domain.User, error) {
	return f.byID, f.err
}

type fakeAudits struct {
	events []domain.AuditEvent
	err    error
}

func (f *fakeAudits) Record(_ context.Context, event domain.AuditEvent) error {
	f.events = append(f.events, event)
	return f.err
}

type fakePasswords struct{ valid bool }

func (f fakePasswords) Verify(string, string) error {
	if !f.valid {
		return errors.New("mismatch")
	}
	return nil
}

type fakeTokens struct{ err error }

func (f fakeTokens) Issue(string, domain.Role) (string, error) {
	return "signed-token", f.err
}

func activeUser(role domain.Role) domain.User {
	return domain.User{
		ID: "user-1", EmployeeID: "12345678", DisplayName: "Example",
		PasswordHash: "hash", Role: role, Active: true,
	}
}

func TestLoginSuccessAuditsAndReturnsToken(t *testing.T) {
	users := &fakeUsers{byEmployee: activeUser(domain.RoleViewer)}
	audits := &fakeAudits{}
	uc := NewLogin(users, audits, fakePasswords{valid: true}, fakeTokens{}, "dummy")

	output, err := uc.Execute(context.Background(), LoginInput{
		EmployeeID: "12345678", Password: "secret", RequestID: "request-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if output.Token != "signed-token" || output.User.ID != "user-1" {
		t.Fatalf("unexpected output: %#v", output)
	}
	if len(audits.events) != 1 || audits.events[0].Outcome != "success" || audits.events[0].RequestID != "request-1" {
		t.Fatalf("unexpected audits: %#v", audits.events)
	}
}

func TestLoginFailuresAreSafeAndAudited(t *testing.T) {
	tests := []struct {
		name string
		user domain.User
		err  error
		pass bool
		want domain.ErrorKind
	}{
		{name: "unknown employee", err: domain.NewError(domain.KindNotFound, nil), want: domain.KindInvalidCredentials},
		{name: "wrong password", user: activeUser(domain.RoleViewer), want: domain.KindInvalidCredentials},
		{name: "inactive", user: func() domain.User { u := activeUser(domain.RoleViewer); u.Active = false; return u }(), pass: true, want: domain.KindAccountInactive},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			audits := &fakeAudits{}
			uc := NewLogin(&fakeUsers{byEmployee: test.user, err: test.err}, audits, fakePasswords{valid: test.pass}, fakeTokens{}, "dummy")
			_, err := uc.Execute(context.Background(), LoginInput{EmployeeID: "x", Password: "secret", RequestID: "r"})
			if domain.ErrorKindOf(err) != test.want {
				t.Fatalf("error kind = %q, want %q", domain.ErrorKindOf(err), test.want)
			}
			if len(audits.events) != 1 || audits.events[0].Outcome != "failure" {
				t.Fatalf("failure not audited: %#v", audits.events)
			}
		})
	}
}

func TestLoginFailsClosedWhenAuditFails(t *testing.T) {
	uc := NewLogin(
		&fakeUsers{byEmployee: activeUser(domain.RoleAdmin)},
		&fakeAudits{err: errors.New("database unavailable")},
		fakePasswords{valid: true},
		fakeTokens{},
		"dummy",
	)
	if _, err := uc.Execute(context.Background(), LoginInput{RequestID: "r"}); domain.ErrorKindOf(err) != domain.KindInternal {
		t.Fatalf("error = %v", err)
	}
}

func TestCurrentUserEnforcesReadRoleAndReloadsAccount(t *testing.T) {
	for _, role := range []domain.Role{domain.RoleViewer, domain.RoleAdmin} {
		user := activeUser(role)
		uc := NewCurrentUser(&fakeUsers{byID: user})
		got, err := uc.Execute(context.Background(), domain.Principal{UserID: user.ID, Role: role})
		if err != nil || got.ID != user.ID {
			t.Fatalf("role %q rejected: user=%#v err=%v", role, got, err)
		}
	}

	uc := NewCurrentUser(&fakeUsers{byID: activeUser(domain.RoleViewer)})
	if _, err := uc.Execute(context.Background(), domain.Principal{UserID: "user-1", Role: "operator"}); domain.ErrorKindOf(err) != domain.KindForbidden {
		t.Fatalf("unexpected role error: %v", err)
	}

	inactive := activeUser(domain.RoleViewer)
	inactive.Active = false
	uc = NewCurrentUser(&fakeUsers{byID: inactive})
	if _, err := uc.Execute(context.Background(), domain.Principal{UserID: "user-1", Role: domain.RoleViewer}); domain.ErrorKindOf(err) != domain.KindAccountInactive {
		t.Fatalf("inactive error: %v", err)
	}
}

func TestLogoutAuditsOnlyResolvedPrincipal(t *testing.T) {
	audits := &fakeAudits{}
	uc := NewLogout(audits)
	if err := uc.Execute(context.Background(), nil, "request-1"); err != nil {
		t.Fatal(err)
	}
	principal := &domain.Principal{UserID: "user-1", Role: domain.RoleViewer}
	if err := uc.Execute(context.Background(), principal, "request-2"); err != nil {
		t.Fatal(err)
	}
	if len(audits.events) != 1 || audits.events[0].Action != "logout" || audits.events[0].ActorUserID == nil {
		t.Fatalf("unexpected audits: %#v", audits.events)
	}
}
