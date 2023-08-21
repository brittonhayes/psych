package sqlite

import (
	"context"

	"github.com/brittonhayes/therapy/api"
	"github.com/uptrace/bun"
)

func (r *repository) therapistFilterQuery(query *bun.SelectQuery, params *api.GetTherapistParams) (*bun.SelectQuery, error) {
	if params == nil {
		return query, nil
	}

	if params.Limit != nil {
		query = query.Limit(*params.Limit)
	}

	if params.Offset != nil {
		query = query.Offset(*params.Offset)
	}

	if params.Title != nil {
		query.Where("? LIKE ?", bun.Ident("title"), *params.Title+"%")
	}

	if params.Credentials != nil {
		query.Where("? LIKE ?", bun.Ident("credentials"), "%"+*params.Credentials+"%")
	}

	if params.Verified != nil {
		query.Where("? LIKE ?", bun.Ident("verified"), "%"+*params.Verified+"%")
	}

	if params.Statement != nil {
		query.Where("? LIKE ?", bun.Ident("statement"), "%"+*params.Statement+"%")
	}

	if params.Phone != nil {
		query.Where("? LIKE ?", bun.Ident("phone"), "%"+*params.Phone+"%")
	}

	if params.Location != nil {
		query.Where("? LIKE ?", bun.Ident("location"), "%"+*params.Location+"%")
	}

	return query, nil
}

func (r *repository) Save(ctx context.Context, therapist api.Therapist) error {
	_, err := r.db.NewInsert().Model(&therapist).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Find(ctx context.Context, params *api.GetTherapistParams) ([]api.Therapist, error) {
	var therapists []api.Therapist

	query, err := r.therapistFilterQuery(r.db.NewSelect().Model(&therapists), params)
	if err != nil {
		return nil, err
	}

	err = query.Scan(ctx, &therapists)
	if err != nil {
		return nil, err
	}

	return therapists, nil
}

func (r *repository) List(ctx context.Context) ([]api.Therapist, error) {
	var therapists []api.Therapist
	err := r.db.NewSelect().Model(&therapists).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return therapists, nil
}
