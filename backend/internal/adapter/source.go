package adapter

import (
	"context"

	"github.com/mohammadalviyan/rpmp/backend/internal/domain"
)

// Source is the anti-corruption boundary for operational data.
// Implementations map vendor payloads to RPMP domain values and never expose
// vendor types to handlers or the API contract.
type Source interface {
	Snapshot(context.Context) (domain.SourceSnapshot, error)
}
