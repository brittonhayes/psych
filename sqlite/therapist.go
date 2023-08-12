package sqlite

import (
	"context"

	"github.com/brittonhayes/therapy"
)

func (r *repository) Save(ctx context.Context, therapist therapy.Therapist) error {
	_, err := r.db.NewInsert().Model(&therapist).Exec(ctx)
	if err != nil {
		return err
	}

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

func (r *repository) List(ctx context.Context) ([]therapy.Therapist, error) {
	var therapists []therapy.Therapist
	err := r.db.NewSelect().Model(&therapists).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return therapists, nil
}
