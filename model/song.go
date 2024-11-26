package model

import (
	"errors"
	"time"
)

type SongFilter struct {
	Performer int       `json:"performer"`
	Song      string    `json:"song"`
	Lyric     string    `json:"lyric"`
	Link      string    `json:"link"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type SongInput struct {
	Name      string `json:"name" binding:"required"`
	Performer int    `json:"performer" binding:"required"`
}

type SongOutput struct {
	Id          int       `json:"-" db:"id"`
	Name        string    `json:"song" db:"name" example:"Ooh Baby"`
	ReleaseDate time.Time `json:"release_date" db:"release_date" example:"2020-01-01"`
	Lyric       *string   `json:"lyric" db:"lyric" example:"Ooh baby, don't you know I suffer?\\nOoh baby, can you hear me moan?"`
	Link        *string   `json:"link" db:"link" example:"https://www.youtube.com/"`
	Performer   string    `json:"performer" db:"performer" example:"singer"`
}

type PaginatedSong struct {
	Songs  []SongOutput `json:"songs"`
	LastId int          `json:"last_id"`
}

type PaginatedLyric struct {
	Lyrics []string `json:"lyrics"`
}

type UpdateSongInput struct {
	Name        *string `json:"name"`
	ReleaseDate *string `json:"release_date"`
	Lyric       *string `json:"lyric"`
	Link        *string `json:"link"`
	Performer   *string `json:"performer_id"`
}

func (u *UpdateSongInput) Validate() error {
	if u.Name == nil && u.Performer == nil && u.Link == nil && u.ReleaseDate == nil && u.Lyric == nil {
		return errors.New("nothing to update")
	}
	return nil
}
