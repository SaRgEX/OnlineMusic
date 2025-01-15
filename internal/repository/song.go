package repository

import (
	"OnlineMusic/model"
	"OnlineMusic/pkg/logger"
	"OnlineMusic/utils"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type SongRepository struct {
	conn   *pgx.Conn
	logger *logger.Logger
	qb     *utils.QueryBuilder
}

func NewSongRepository(conn *pgx.Conn, logger *logger.Logger, qb *utils.QueryBuilder) *SongRepository {
	return &SongRepository{conn: conn, logger: logger, qb: qb}
}

func (r *SongRepository) ViewAll(ctx context.Context, input model.SongFilter, cursor, pageSize int) ([]model.SongOutput, error) {
	query := fmt.Sprintf(`SELECT id, name, lyric, release_date, link, performer 
		FROM %s($1, $2, $3, $4, $5, $6, $7, $8) AS s`, songFilter)
	rows, err := r.conn.Query(ctx, query, input.Song, input.Performer, input.Lyric, input.Link, input.StartDate, input.EndDate, cursor, pageSize)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	r.logger.With("method", "ViewAll", "query", query).Info("Viewing all songs")
	songs, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.SongOutput])
	return songs, err
}

func (r *SongRepository) FindText(ctx context.Context, id int) (*string, error) {
	var lyric *string
	query := fmt.Sprintf(`SELECT lyric FROM %s WHERE id = $1`, "song")
	err := r.conn.QueryRow(ctx, query, id).Scan(&lyric)
	r.logger.With("method", "FindText", "query", query).Info("Find song's lyrics")
	return lyric, err
}

func (r *SongRepository) Add(ctx context.Context, song model.SongInput) error {
	query := fmt.Sprintf(`INSERT INTO %s (name, performer_id) VALUES ($1, $2)`, songTable)
	_, err := r.conn.Exec(ctx, query, song.Name, song.Performer)
	r.logger.With("method", "Add", "query", query).Info("Adding song")
	return err
}

func (r *SongRepository) Delete(ctx context.Context, id int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, songTable)
	_, err := r.conn.Exec(ctx, query, id)
	r.logger.With("method", "Delete", "query", query).Info("Delete song")
	return err
}

func (r *SongRepository) Update(ctx context.Context, id int, song model.UpdateSongInput) error {
	query, args := r.qb.BuildUpdateQueryFromSong(songTable, "id", id, song)
	_, err := r.conn.Exec(ctx, query, args...)
	r.logger.With("method", "Update", "query", query).Info("Update song")
	return err
}
