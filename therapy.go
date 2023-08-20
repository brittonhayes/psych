//go:generate go run github.com/99designs/gqlgen generate
package therapy

import (
	"context"

	"github.com/brittonhayes/therapy/api"
)

type Repository interface {
	Save(ctx context.Context, therapist api.Therapist) error
	Find(ctx context.Context, therapist *api.GetTherapistParams) ([]api.Therapist, error)
	List(ctx context.Context) ([]api.Therapist, error)

	Init(ctx context.Context) error
	Generate(ctx context.Context, name string) error
	Migrate(ctx context.Context) error
	Lock(ctx context.Context) error
	Unlock(ctx context.Context) error
	Rollback(ctx context.Context) error
}
