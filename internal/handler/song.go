package handler

import (
	"OnlineMusic/model"
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

// @Summary Get all songs
// @Description Get all songs with filtration by each column and pagination
// @Tags songs
// @Accept  json
// @Produce  json
// @Param performer_id query int false "Performer"
// @Param song query string false "Song"
// @Param startDate query string false "Start Date" Format(date)
// @Param endDate query string false "End Date" Format(date)
// @Param link query string false "Link"
// @Param lyric query string false "Lyrics"
// @Param cursor query int false "Cursor" default("1")
// @Param pageSize query int false "Page size" default("10")
// @Success 200 {object} model.PaginatedSong "Songs"
// @Failure 400 {object} errorResponse "Incorrect input"
// @Failure 404 {object} errorResponse "Not found"
// @Failure 500 {object} errorResponse "Server error"
// @Router /songs [get]
func (h *Handler) viewAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	performer, err := GetIntQueryParam(c, "performer_id", -1)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	cursor, err := GetIntQueryParam(c, "cursor", 0)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	pageSize, err := GetIntQueryParam(c, "pageSize", 10)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	startDate, err := ParseDateParam(c, "startDate", "1970-01-01")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	endDate, err := ParseDateParam(c, "endDate", "9999-01-01")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	input := model.SongFilter{
		Song:      c.Query("song"),
		Performer: performer,
		Link:      c.Query("link"),
		Lyric:     c.Query("lyric"),
		StartDate: startDate,
		EndDate:   endDate,
	}
	paginatedSong, err := h.s.ViewAll(ctx, input, cursor, pageSize)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if len(paginatedSong.Songs) == 0 {
		newErrorResponse(c, http.StatusNotFound, "Not found")
		return
	}
	c.JSON(http.StatusOK, paginatedSong)
}

// @Summary Get song's text with pagination
// @Description Get song by id
// @Tags songs
// @Produce  json
// @Param id path int true "ID"
// @Param cursor query int false "Cursor" default("1")
// @Param pageSize query int false "Page size" default("10")
// @Success 200 {object} model.PaginatedLyric "Lyrics"
// @Failure 400 {object} errorResponse "Incorrect input"
// @Failure 404 {object} errorResponse "Not found"
// @Failure 500 {object} errorResponse "Server error"
// @Router /songs/{id}/lyrics [get]
func (h *Handler) findText(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	cursor, err := strconv.Atoi(c.Query("cursor"))
	if err != nil || cursor <= 0 {
		cursor = 1
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}
	lyric, err := h.s.FindText(ctx, id, cursor-1, pageSize)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if len(lyric.Lyrics) == 0 {
		newErrorResponse(c, http.StatusBadRequest, "There are no lyrics")
		return
	}
	c.JSON(http.StatusOK, lyric)
}

// @Summary Add song
// @Description Add song by model.SongInput. Make a fetch request to the third party api.
// @Tags songs
// @Accept  json
// @Produce  json
// @Param song body model.SongInput true "Song"
// @Success 200 {object} statusResponse "Status"
// @Failure 400 {object} errorResponse "Incorrect input"
// @Failure 500 {object} errorResponse "Server error"
// @Router /songs/ [post]
func (h *Handler) add(c *gin.Context) {
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var input model.SongInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.s.Add(dbCtx, input)
	if err != nil {
		slog.With("error", err, "model", input).Error("Error adding song")
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	apiCtx, apiCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer apiCancel()

	resultChan := make(chan model.SongInfoDetail)
	errorChan := make(chan error)

	go func() {
		defer close(resultChan)
		defer close(errorChan)
		info, err := h.c.FetchInfoMusic(apiCtx, input.Name, strconv.Itoa(input.Performer))
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- info
	}()

	select {
	case <-apiCtx.Done():
		slog.Debug("Music info fetched successfully")
		c.JSON(http.StatusCreated, gin.H{
			"info": "Error while fetching music info",
		})
		return
	case err := <-errorChan:
		slog.With("error", err).Error("Error while fetching music info")
		c.JSON(http.StatusCreated, gin.H{
			"info": "Error while fetching music info",
		})
		return
	case info := <-resultChan:
		c.JSON(http.StatusCreated, info)
	}
}

// @Summary Delete song
// @Description Delete song by id
// @Tags songs
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} statusResponse "Status"
// @Failure 400 {object} errorResponse "Incorrect input"
// @Failure 500 {object} errorResponse "Server error"
// @Router /songs/{id} [delete]
func (h *Handler) delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.s.Delete(ctx, id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	c.JSON(http.StatusOK, statusResponse{
		Status: "deleted",
	})
}

// @Summary Update song
// @Description Update song by id from query
// @Tags songs
// @Produce  json
// @Param id path int true "ID"
// @Param song body model.UpdateSongInput true "Song"
// @Success 200 {object} statusResponse "Status"
// @Failure 400 {object} errorResponse "Incorrect input"
// @Failure 500 {object} errorResponse "Server error"
// @Router /songs/{id} [put]
func (h *Handler) update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var input model.UpdateSongInput
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = c.ShouldBindJSON(&input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.s.Update(ctx, id, input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{
		Status: "updated",
	})

}
