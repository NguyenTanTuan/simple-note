package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/nguyentantuan/simple-note/graph/generated"
	"github.com/nguyentantuan/simple-note/graph/model"
	"math/rand"
)

func (r *mutationResolver) AddNote(ctx context.Context, text string, description string) (*model.Note, error) {
	note := &model.Note{
		Text:        text,
		ID:          fmt.Sprintf("T%d", rand.Int()),
		Description: description,
	}
	r.notes = append(r.notes, note)
	return note, nil
}

func (r *mutationResolver) AddSubNote(ctx context.Context, text string, description string, idParent string) (*model.SubNote, error) {
	var note *model.Note
	for _, n := range r.notes {
		if n.ID == idParent {
			note = n
			break
		}
	}

	if note == nil {
		panic("Parent note not found")
	}

	newSubNote := &model.SubNote{
		ID:          fmt.Sprintf("T%d", rand.Int()),
		Text:        text,
		Description: description,
	}

	note.SubNote = append(note.SubNote, newSubNote)
	return newSubNote, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
