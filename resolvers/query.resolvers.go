package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/nguyentantuan/simple-note/graph/generated"
	"github.com/nguyentantuan/simple-note/graph/model"
)

func (r *queryResolver) GetAllNotes(ctx context.Context) ([]*model.Note, error) {
	return r.notes, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
