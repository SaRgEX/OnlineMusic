package utils

import (
	"OnlineMusic/model"
	"fmt"
	"strings"
)

type BuildQuery interface {
	BuildUpdateQueryFromSong(table string, idColumn string, idValue interface{}, input model.UpdateSongInput) (string, []interface{})
	BuildUpdateQuery(table string, idColumn string, idValue interface{}, fields map[string]interface{}) (string, []interface{})
}

type QueryBuilder struct {
	BuildQuery
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

func (qb *QueryBuilder) BuildUpdateQueryFromSong(table string, idColumn string, idValue interface{}, input model.UpdateSongInput) (string, []interface{}) {
	fields := make(map[string]interface{})

	if input.Name != nil {
		fields["name"] = *input.Name
	}
	if input.ReleaseDate != nil {
		fields["release_date"] = *input.ReleaseDate
	}
	if input.Lyric != nil {
		fields["lyric"] = *input.Lyric
	}
	if input.Link != nil {
		fields["link"] = *input.Link
	}
	if input.Performer != nil {
		fields["creator"] = *input.Performer
	}

	return qb.BuildUpdateQuery(table, idColumn, idValue, fields)
}

func (qb *QueryBuilder) BuildUpdateQuery(table string, idColumn string, idValue interface{}, fields map[string]interface{}) (string, []interface{}) {
	setParts := make([]string, 0, len(fields))
	args := make([]interface{}, 0, len(fields)+1)
	argID := len(fields) + 1

	i := 1
	for column, value := range fields {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", column, i))
		args = append(args, value)
		i++
	}

	args = append(args, idValue)

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s = $%d",
		table,
		strings.Join(setParts, ", "),
		idColumn,
		argID,
	)

	return query, args
}
