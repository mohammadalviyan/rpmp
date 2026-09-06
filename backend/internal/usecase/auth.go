package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mohammadalviyan/rpmp/backend/internal/domain"
)

type UserRepository interface {
	FindByEmployeeID(context.Context, string) (domain.User, error)
	FindByID(context.Context, string) (domain.User, error)
}

type AuditRepository interface {
	Record(context.Context, domain.AuditEvent) error
}

type PasswordVerifier interface {
	Verify(hash, password string) error
}

type TokenIssuer interface {
	Issue(userID string, role domain.Role) (string, error)
}

type LoginInput struct {
	EmployeeID string
	Password   string
	RequestID  string
}

type LoginOutput struct {
	User  domain.User
	Token string
}

type Login struct {
	users             UserRepository
	audits            AuditRepository
	passwords         PasswordVerifier
	tokens            TokenIssuer
	dummyPasswordHash string
	now               func() time.Time
	newID             func() string
}

func NewLogin(users UserRepository, audits AuditRepository, passwords PasswordVerifier, tokens TokenIssuer, dummyPasswordHash string) *Login {
	return &Login{
		users: users, audits: audits, passwords: passwords, tokens: tokens,
		dummyPasswordHash: dummyPasswordHash,
		now:               time.Now,
		newID:             func() string { return uuid.NewString() },
	}
}

func (uc *Login) Execute(ctx context.Context, input LoginInput) (LoginOutput, error) {
	user, err := uc.users.FindByEmployeeID(ctx, input.EmployeeID)
	if err != nil {
		if domain.ErrorKindOf(err) != domain.KindNotFound {
			return LoginOutput{}, domain.NewError(domain.KindInternal, err)
		}
		_ = uc.passwords.Verify(uc.dummyPasswordHash, input.Password)
		if auditErr := uc.record(ctx, nil, "login", "failure", input.RequestID); auditErr != nil {
			return LoginOutput{}, domain.NewError(domain.KindInternal, auditErr)
		}
		return LoginOutput{}, domain.NewError(domain.KindInvalidCredentials, nil)
	}

	actorID := &user.ID
	if err := uc.passwords.Verify(user.PasswordHash, input.Password); err != nil {
		if auditErr := uc.record(ctx, actorID, "login", "failure", input.RequestID); auditErr != nil {
			return LoginOutput{}, domain.NewError(domain.KindInternal, auditErr)
		}
		return LoginOutput{}, domain.NewError(domain.KindInvalidCredentials, nil)
	}
	if !user.Active {
		if auditErr := uc.record(ctx, actorID, "login", "failure", input.RequestID); auditErr != nil {
			return LoginOutput{}, domain.NewError(domain.KindInternal, auditErr)
		}
		return LoginOutput{}, domain.NewError(domain.KindAccountInactive, nil)
	}
	if !user.Role.CanRead() {
		if auditErr := uc.record(ctx, actorID, "login", "failure", input.RequestID); auditErr != nil {
			return LoginOutput{}, domain.NewError(domain.KindInternal, auditErr)
		}
		return LoginOutput{}, domain.NewError(domain.KindForbidden, nil)
	}

	token, err := uc.tokens.Issue(user.ID, user.Role)
	if err != nil {
		if auditErr := uc.record(ctx, actorID, "login", "failure", input.RequestID); auditErr != nil {
			return LoginOutput{}, domain.NewError(domain.KindInternal, auditErr)
		}
		return LoginOutput{}, domain.NewError(domain.KindInternal, err)
	}
	if err := uc.record(ctx, actorID, "login", "success", input.RequestID); err != nil {
		return LoginOutput{}, domain.NewError(domain.KindInternal, err)
	}
	return LoginOutput{User: user, Token: token}, nil
}

func (uc *Login) record(ctx context.Context, actorID *string, action, outcome, requestID string) error {
	return uc.audits.Record(ctx, domain.AuditEvent{
		ID: uc.newID(), ActorUserID: actorID, Action: action, Outcome: outcome,
		RequestID: requestID, OccurredAt: uc.now().UTC(),
	})
}

type CurrentUser struct {
	users UserRepository
}

func NewCurrentUser(users UserRepository) *CurrentUser {
	return &CurrentUser{users: users}
}

func (uc *CurrentUser) Execute(ctx context.Context, principal domain.Principal) (domain.User, error) {
	if !principal.Role.CanRead() {
		return domain.User{}, domain.NewError(domain.KindForbidden, nil)
	}
	user, err := uc.users.FindByID(ctx, principal.UserID)
	if err != nil {
		if domain.ErrorKindOf(err) == domain.KindNotFound {
			return domain.User{}, domain.NewError(domain.KindUnauthenticated, err)
		}
		return domain.User{}, domain.NewError(domain.KindInternal, err)
	}
	if !user.Active {
		return domain.User{}, domain.NewError(domain.KindAccountInactive, nil)
	}
	if !user.Role.CanRead() {
		return domain.User{}, domain.NewError(domain.KindForbidden, nil)
	}
	return user, nil
}

type Logout struct {
	audits AuditRepository
	now    func() time.Time
	newID  func() string
}

func NewLogout(audits AuditRepository) *Logout {
	return &Logout{
		audits: audits,
		now:    time.Now,
		newID:  func() string { return uuid.NewString() },
	}
}

func (uc *Logout) Execute(ctx context.Context, principal *domain.Principal, requestID string) error {
	if principal == nil {
		return nil
	}
	actorID := principal.UserID
	if actorID == "" {
		return nil
	}
	err := uc.audits.Record(ctx, domain.AuditEvent{
		ID: uc.newID(), ActorUserID: &actorID, Action: "logout", Outcome: "success",
		RequestID: requestID, OccurredAt: uc.now().UTC(),
	})
	if err != nil {
		return domain.NewError(domain.KindInternal, err)
	}
	return nil
}
