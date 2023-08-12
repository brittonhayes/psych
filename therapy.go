package therapy

import (
	"context"
)

type Therapist struct {
	ID          int    `bun:"id,pk,autoincrement" json:"id"`
	Title       string `json:"title"`
	Credentials string `json:"credentials"`
	Verified    string `json:"verified"`
	Statement   string `json:"statement"`
	Phone       string `json:"phone"`
	Location    string `json:"location"`
	Link        string `bun:",unique" json:"link"`
}

type Repository interface {
	Save(ctx context.Context, therapist Therapist) error
	Find(ctx context.Context, therapist Therapist) ([]Therapist, error)
	List(ctx context.Context) ([]Therapist, error)

	Init(ctx context.Context) error
	Generate(ctx context.Context, name string) error
	Migrate(ctx context.Context) error
	Lock(ctx context.Context) error
	Unlock(ctx context.Context) error
	Rollback(ctx context.Context) error
}
