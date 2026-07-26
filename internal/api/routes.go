package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/raffleberry/riptvtime/internal/db"
	"github.com/raffleberry/riptvtime/internal/services"
)

var (
	ErrInvalidRequest = errors.New("Invalid Request")
)

// queryParams = { q: query, p: page }
func (a *Api) SeriesSearch() http.HandlerFunc {
	return WithCtx(func(c *Context) error {
		urlVals := c.R.URL.Query()
		queryStr := urlVals.Get("q")
		pageStr := urlVals.Get("p")
		slog.Debug("search", "page", pageStr, "query", queryStr)

		var err error
		page := 0

		if pageStr != "" {
			page, err = strconv.Atoi(pageStr)
		}
		if err != nil {
			return c.Error(http.StatusBadRequest, err.Error())
		}

		respRes, err := a.tv.Search(queryStr, page)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, respRes)
	})
}

func (a *Api) SeriesAdd() http.HandlerFunc {
	return WithCtx(func(c *Context) error {
		var payload struct {
			MId int
		}

		if err := json.NewDecoder(c.R.Body).Decode(&payload); err != nil {
			return err
		}

		slog.Debug("tv add", "payload", payload)

		s, err := a.tv.Add(payload.MId)

		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, s)
	})
}

func (a *Api) SeriesEpisodeWatch() http.HandlerFunc {
	return WithCtx(func(ctx *Context) error {
		var payload struct {
			SeriesMId int
			SeasonNo  int
			EpisodeNo int
		}

		if err := json.NewDecoder(ctx.R.Body).Decode(&payload); err != nil {
			return err
		}

		insertId, err := a.tv.SetEpisodeWatched(payload.SeriesMId, payload.SeasonNo, payload.EpisodeNo)
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, insertId)
	})
}

func (a *Api) SeriesEpisodeUnWatch() http.HandlerFunc {
	return WithCtx(func(ctx *Context) error {
		var payload struct {
			SeriesMId int
			SeasonNo  int
			EpisodeNo int
		}

		if err := json.NewDecoder(ctx.R.Body).Decode(&payload); err != nil {
			return err
		}

		err := a.tv.SetEpisodeUnwatch(payload.SeriesMId, payload.SeasonNo, payload.EpisodeNo)
		if err != nil {
			if errors.Is(err, services.ErrNotFound) {
				return ctx.Error(http.StatusNotFound, err.Error())
			}
			return err
		}
		return ctx.JSON(http.StatusOK, nil)
	})
}

func (a *Api) SeriesFeed() http.HandlerFunc {
	return WithCtx(func(c *Context) error {
		respItem, err := a.tv.Feed()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, respItem)
	})
}

func (a *Api) SeriesUpdateStatus() http.HandlerFunc {
	return WithCtx(func(c *Context) error {
		mIdStr := c.R.PathValue("mId")

		mId, err := strconv.Atoi(mIdStr)

		if err != nil {
			return c.Error(http.StatusBadRequest, err.Error())
		}

		var payload struct {
			Status int
		}

		if err := json.NewDecoder(c.R.Body).Decode(&payload); err != nil {
			return err
		}
		status := db.TvStatus(payload.Status)
		if !status.IsValid() {
			return c.Error(http.StatusBadRequest, "Invalid status")
		}

		err = a.tv.UpdateStatus(mId, status)
		if err != nil {
			if errors.Is(err, services.ErrNotFound) {
				return c.Error(http.StatusNotFound, err.Error())
			}
		}

		return nil
	})
}

func (a *Api) SeriesRem() http.HandlerFunc {
	return WithCtx(func(c *Context) error {
		mIdStr := c.R.PathValue("mId")

		mId, err := strconv.Atoi(mIdStr)

		if err != nil {
			return c.Error(http.StatusBadRequest, err.Error())
		}

		err = a.tv.Remove(mId)

		if err != nil {
			if errors.Is(err, services.ErrNotFound) {
				return c.Error(http.StatusNotFound, err.Error())
			}
		}

		return nil

	})
}

func (a *Api) SeriesGet() http.HandlerFunc {
	return WithCtx(func(c *Context) error {
		mIdStr := c.R.PathValue("mId")
		full := c.R.URL.Query().Get("full")
		var sd bool = len(full) > 0

		mId, err := strconv.Atoi(mIdStr)

		if err != nil {
			return c.Error(http.StatusBadRequest, err.Error())
		}
		rv, err := a.tv.GetDetails(mId, sd)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, rv)
	})
}

func (a *Api) SeriesUpNext() http.HandlerFunc {
	return WithCtx(func(c *Context) error {
		mIdStr := c.R.PathValue("mId")
		mId, err := strconv.Atoi(mIdStr)

		if err != nil {
			return c.Error(http.StatusBadRequest, err.Error())
		}
		rv, err := a.tv.UpNext(mId)

		if errors.Is(err, services.ErrNotFound) {
			return c.Text(http.StatusNoContent, "")
		} else if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, rv)
	})
}

func (a *Api) SeriesAll() http.HandlerFunc {
	return WithCtx(func(c *Context) error {
		respItem, err := a.tv.TrackedAll()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, respItem)
	})
}
