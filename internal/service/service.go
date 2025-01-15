package service

import (
	"OnlineMusic/internal/repository"
	"OnlineMusic/model"
	"OnlineMusic/pkg/logger"
	"context"
)

type Song interface {
	ViewAll(ctx context.Context, input model.SongFilter, cursor int, pageSize int) (model.PaginatedSong, error)
	FindText(ctx context.Context, id, cursor, pageSize int) (model.PaginatedLyric, error)
	Add(ctx context.Context, input model.SongInput) error
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, input model.UpdateSongInput) error
}

type Service struct {
	Song
}

func New(r *repository.Repository, logger *logger.Logger) *Service {
	return &Service{
		Song: NewSongService(r.Song, logger),
	}
}
