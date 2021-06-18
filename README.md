# simple-note
This a separate section, We will build another simple note from scratch
Create database and generate your repositories by
Create new migration
Migrate file name should has increase prefix

migrations/00068_add_note.sql

-- +goose Up

CREATE TABLE IF NOT EXISTS `note` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `text` varchar(255) NOT NULL,
    `description` TEXT,
    PRIMARY KEY (`id`)
) ENGINE=INNODB;

CREATE TABLE IF NOT EXISTS `sub_note` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `fk_note` INT NOT NULL,
    `text` varchar(255) NOT NULL,
    `description` TEXT,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`fk_note`) REFERENCES `note`(`id`) ON DELETE NO ACTION
) ENGINE=INNODB;
Clean migrations
$ ./do.sh clean-migrate


Install goimport
go get github.com/polaris1119/go.tools/cmd/goimports
Generate entities, repo
previous part, we generate entities by gqlgen but if you already have a table on database, you can generate entities, graphql type and repositories from your tabe:

$ ./do.sh gen-repo
All generated files by gen-repo will haves xo suffix
entities/note.xo.go
entities/sub_note.xo.go
repositories/note_repository.xo.go
repositories/note_rlts_repository.xo.go
repositories/sub_note_repository.xo.go
repositories/sub_note_rlts_repository.xo.go


// services/note_service.go 

package services

import (
    "context"
    "database/sql"
    "github.com/google/wire"
    "github.com/pkg/errors"

    "sn-backend/entities"
    "sn-backend/repositories"
)

type INoteService interface {
    AddNote(ctx context.Context, text string, description *string) (*entities.Note, error)
    GetAllNotes(ctx context.Context, filter *entities.NoteFilter, pagination *entities.Pagination) (entities.ListNote, error)
    AddSubNote(ctx context.Context, IDParent int, text string, description *string) (*entities.SubNote, error)
}

type NoteServiceOptions struct {
    NoteRepository      repositories.NoteRepository
    SubNoteRepository   repositories.SubNoteRepository
}

type NoteService struct {
    NoteRepository      repositories.NoteRepository
    SubNoteRepository   repositories.SubNoteRepository
}

var NewNoteService = wire.NewSet(wire.Struct(new(NoteService), "*"), wire.Bind(new(INoteService), new(NoteService)))

func (vs *NoteService) AddNote(ctx context.Context, text string, description *string) (*entities.Note, error) {
    desc, _ := entities.UnmarshalNullString(description)
    noteCreate := entities.NoteCreate{
        Text: text,
        Description: desc,
    }

    return vs.NoteRepository.InsertNote(ctx, noteCreate)
}

func (vs *NoteService) GetAllNotes(ctx context.Context, filter *entities.NoteFilter, pagination *entities.Pagination) (entities.ListNote, error) {
    return vs.NoteRepository.FindAllNote(ctx, filter, pagination)
}

func (vs *NoteService) AddSubNote(ctx context.Context, IDParent int, text string, description *string) (*entities.SubNote, error) {
    _, err := vs.NoteRepository.NoteByID(ctx, IDParent, nil)
    if err != nil {
        if errors.Cause(err) == sql.ErrNoRows {
            return nil, errors.New("Note not found")
        } else {
            return nil, err
        }
    }

    desc, err := entities.UnmarshalNullString(description)
    newSubNote := entities.SubNoteCreate{
        Text:           text,
        Description:    desc,
        FkNote:         IDParent,
    }

    return vs.SubNoteRepository.InsertSubNote(ctx, newSubNote)
}


Register note service and repositories
NoteService now use repositories so, in resolver.go, let declare new resolver:

// graphql/resolver.go
...
type Resolver struct {
    ...
    NoteRltsRepository      repositories.INoteRltsRepository
    SubNoteRltsRepository   repositories.ISubNoteRltsRepository
}
...
func (r *Resolver) Note() NoteResolver  { return r.NoteRltsRepository }
func (r *Resolver) SubNote() SubNoteResolver  { return r.SubNoteRltsRepository }
...
In the wire.go config file, you need to register generated repositories

// internal/wire.go

...
var repositorySet = wire.NewSet(
    ...
    repositories.NewNoteRepository,
    repositories.NewNoteRltsRepository, 
    repositories.NewSubNoteRepository,
    repositories.NewSubNoteRltsRepository,
)
...
Then wire again
$ ./do.sh wire
Define your mutations within mutation.graphql

# /graphql/schema/mutation.graphql

...
type Mutation {
  addNote(text: String!, description: String): Note
  addSubNote(IDParent: Int!, text: String!, description: String ): SubNote
    ...
}

...
Define your query within query.graphql
# /graphql/schema/query.graphql

...
type Query {
    getAllNotes(filter: NoteFilter,  pagination: Pagination ): ListNote!
    ...
}
...
Generate connecting code by wire
In resolver.go, add INoteService as resolver for note

// graphql/resolver.go 

...
type Resolver struct {
    ...
    services.INoteService

}
In wire.go add NewNoteService as provider to inject
// internal/wire.go

...
var servicesSet = wire.NewSet(
    ...
    services.NewNoteService,
)
...
Generate connecting code
./do.sh wire

./do.sh gen-gql
Here are some queries to try:

#!graphQL

mutation createNote {
  addNote(text:"todo11", description:"11") {
    text
    description
  }
}

query getAllNotes {
  getAllNotes {
    data {
      text
    }
  }
}
Done! The video is attached below



