package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"context"

	"github.com/brittonhayes/therapy"
	"github.com/brittonhayes/therapy/api"
)

// Therapists is the resolver for the therapists field.
func (r *queryResolver) Therapists(ctx context.Context, filter *therapy.TherapistFilters) ([]api.Therapist, error) {
	if filter == nil {
		return r.Repo.List(ctx)
	}

	return r.Repo.Find(ctx, &api.GetTherapistParams{
		Title:       filter.Title,
		Credentials: filter.Credentials,
		Verified:    filter.Verified,
		Statement:   filter.Statement,
		Phone:       filter.Phone,
		Location:    filter.Location,
		Link:        filter.Link,
		Limit:       filter.Limit,
		Offset:      filter.Offset,
	})
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }