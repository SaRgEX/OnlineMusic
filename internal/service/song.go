package service

import (
	"OnlineMusic/internal/repository"
	"OnlineMusic/model"
	"OnlineMusic/pkg/logger"
	"context"
	"strings"
)

type SongService struct {
	r      repository.Song
	logger *logger.Logger
}

func NewSongService(repository repository.Song, logger *logger.Logger) *SongService {
	return &SongService{
		r:      repository,
		logger: logger,
	}
}

func (s *SongService) ViewAll(ctx context.Context, filter model.SongFilter, cursor int, pageSize int) (response model.PaginatedSong, err error) {
	s.logger.With("ctx", ctx)
	response.Songs, err = s.r.ViewAll(ctx, filter, cursor, pageSize)
	if len(response.Songs) <= 0 {
		return response, err
	}
	response.LastId = response.Songs[len(response.Songs)-1].Id
	return response, err
}

func (s *SongService) FindText(ctx context.Context, id, cursor, pageSize int) (model.PaginatedLyric, error) {
	var lyric model.PaginatedLyric
	lyrics, err := s.r.FindText(ctx, id)
	if err != nil || lyrics == nil {
		return model.PaginatedLyric{}, err
	}
	lyric.Lyrics = paginateText(*lyrics, cursor, pageSize)
	return lyric, err
}

func (s *SongService) Add(ctx context.Context, input model.SongInput) error {
	return s.r.Add(ctx, input)
}

func (s *SongService) Delete(ctx context.Context, id int) error {
	return s.r.Delete(ctx, id)
}

func (s *SongService) Update(ctx context.Context, id int, song model.UpdateSongInput) error {
	if err := song.ValidateSongInput(); err != nil {
		return err
	}
	return s.r.Update(ctx, id, song)
}

func paginateText(text string, cursor, pageSize int) []string {
	rows := strings.Split(text, "\n")
	if cursor > len(rows) {
		return nil
	}
	end := cursor + pageSize
	if end > len(rows) {
		end = len(rows)
	}
	return rows[cursor:end]
}
