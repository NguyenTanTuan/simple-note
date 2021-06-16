package services

import (
	"context"
	"fmt"
	"github.com/nguyentantuan/simple-note/graph/model"
	"math/rand"
)

// declare resolvers
type IVisitService interface {
	AddNote(ctx context.Context, text string, description string) (*model.Note, error)
	GetAllNotes(ctx context.Context) ([]*model.Note, error)
	AddSubNote(ctx context.Context, text string, description string, IDParent string) (*model.SubNote, error)
}

// use variable to store data, we will learn about
// database later
type VisitService struct {
	notes []*model.Note
}
// implement resolvers
func (vs *VisitService) AddNote(ctx context.Context, text string, description string) (*model.Note, error) {
	newNote := &model.Note{
		ID:          fmt.Sprintf("T%d", rand.Int()),
		Text:        text,
		Description: description,
	}

	vs.notes = append(vs.notes, newNote)
	return newNote, nil
}

func (vs *VisitService) GetAllNotes(ctx context.Context) ([]*model.Note, error) {
	return vs.notes, nil
}
func (vs *VisitService) AddSubNote(ctx context.Context, text string, description string, IDParent string) (*model.SubNote, error) {
	var note *model.Note
	for _, n := range vs.notes {
		if n.ID == IDParent {
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