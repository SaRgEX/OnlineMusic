package repository

import (
	"OnlineMusic/model"
	"OnlineMusic/utils"
	"context"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	Song
}

type Song interface {
	ViewAll(ctx context.Context, input model.SongFilter, cursor, pageSize int) ([]model.SongOutput, error)
	FindText(ctx context.Context, id int) (*string, error)
	Add(ctx context.Context, song model.SongInput) error
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, song model.UpdateSongInput) error
}

func New(conn *pgx.Conn, qb *utils.QueryBuilder) *Repository {
	return &Repository{
		Song: NewSongRepository(conn, qb),
	}
}
