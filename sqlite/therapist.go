package sqlite

import (
	"context"
	"log/slog"

	"github.com/brittonhayes/therapy"
)

func (r *repository) Save(ctx context.Context, therapist therapy.Therapist) error {
	_, err := r.db.NewInsert().Model(&therapist).Exec(ctx)
	if err != nil {
		return err
	}

	r.logger.DebugContext(ctx, "saved therapist", slog.String("title", therapist.Title))
	return nil
}

func (r *repository) Find(ctx context.Context, therapist therapy.Therapist) ([]therapy.Therapist, error) {
	var therapists []therapy.Therapist
	err := r.db.NewSelect().Model(&therapists).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return therapists, nil
}
