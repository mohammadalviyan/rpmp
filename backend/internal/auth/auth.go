package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mohammadalviyan/rpmp/backend/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type Passwords struct {
	cost int
}

func NewPasswords(cost int) (*Passwords, error) {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return nil, fmt.Errorf("bcrypt cost is outside supported range")
	}
	return &Passwords{cost: cost}, nil
}

func (p *Passwords) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), p.cost)
	return string(hash), err
}

func (p *Passwords) Verify(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

type Claims struct {
	Role domain.Role `json:"role"`
	jwt.RegisteredClaims
}

type Tokens struct {
	issuer string
	key    []byte
	ttl    time.Duration
	now    func() time.Time
	newJTI func() string
}

func NewTokens(issuer, signingKey string, ttl time.Duration) (*Tokens, error) {
	if issuer == "" {
		return nil, errors.New("JWT issuer is required")
	}
	if len(signingKey) < 32 {
		return nil, errors.New("JWT signing key must contain at least 32 bytes")
	}
	if ttl <= 0 {
		return nil, errors.New("JWT TTL must be positive")
	}
	return &Tokens{
		issuer: issuer,
		key:    []byte(signingKey),
		ttl:    ttl,
		now:    time.Now,
		newJTI: func() string { return uuid.NewString() },
	}, nil
}

func (t *Tokens) Issue(userID string, role domain.Role) (string, error) {
	if userID == "" || !role.CanRead() {
		return "", domain.NewError(domain.KindInternal, errors.New("invalid token subject or role"))
	}
	now := t.now().UTC()
	claims := Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    t.issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(t.ttl)),
			ID:        t.newJTI(),
		},
	}
	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(t.key)
	if err != nil {
		return "", domain.NewError(domain.KindInternal, err)
	}
	return signed, nil
}

func (t *Tokens) Verify(raw string) (domain.Principal, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		raw,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("unexpected signing method")
			}
			return t.key, nil
		},
		jwt.WithIssuer(t.issuer),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithTimeFunc(t.now),
	)
	if err != nil || !token.Valid || claims.Subject == "" || claims.ID == "" || claims.IssuedAt == nil || !claims.Role.CanRead() {
		return domain.Principal{}, domain.NewError(domain.KindUnauthenticated, err)
	}
	return domain.Principal{UserID: claims.Subject, Role: claims.Role}, nil
}
