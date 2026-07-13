package meta

import (
	"log/slog"
	"strconv"
	"time"

	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/hashicorp/go-retryablehttp"
	"gitlab.com/raffleberry/riptvtime/internal/config"
)

type MetaTmdb struct {
	c *tmdb.Client
}

func (t *MetaTmdb) Name() string {
	return "tmdb"
}

func parseYear(date string) int {
	if len(date) >= 4 {
		var err error
		year, err := strconv.Atoi(date[:4])
		if err != nil {
			slog.Warn("Error parsing year", "date", date, "err", err)
			year = 0
		}
		return year
	}
	slog.Warn("Failed to parse year", "date", date)
	return 0
}

func parseAirDate(date string) time.Time {
	if len(date) == 10 {
		firstAirDate, err := time.Parse(time.DateOnly, date)
		if err != nil {
			slog.Warn("Error parsing AirDate", "date", date, "err", err)
		}
		return firstAirDate
	}
	slog.Warn("Failed to parse date", "date", date)
	return time.Time{}
}

func (t *MetaTmdb) Search(query string, page int) (*TvSearchResults, error) {
	var tvSearchResults TvSearchResults

	opts := map[string]string{}
	if page > 0 {
		opts["page"] = strconv.Itoa(page)
	}

	res, err := t.c.GetSearchTVShow(query, opts)

	if err != nil {
		return nil, err
	}

	tvSearchResults.Page = int(res.Page)
	tvSearchResults.TotalPages = int(res.TotalPages)
	tvSearchResults.TotalResults = int(res.TotalResults)

	for _, v := range res.Results {
		tvSearchResults.Results = append(tvSearchResults.Results, TvSearchResult{
			Id:       int(v.ID),
			Name:     v.Name,
			Overview: v.Overview,
			Year:     parseYear(v.FirstAirDate),
			MName:    t.Name(),
		})
	}

	return &tvSearchResults, nil

}

func (t *MetaTmdb) GetTvDetails(tmdbId int) (*TvDetails, error) {
	res, err := t.c.GetTVDetails(tmdbId, nil)

	seasons := []TvSeason{}

	for _, s := range res.Seasons {
		seasons = append(seasons, TvSeason{
			Id:           int(s.ID),
			Name:         s.Name,
			Overview:     s.Overview,
			SeasonNumber: int(s.SeasonNumber),
			EpisodeCount: int(s.EpisodeCount),
			AirDate:      parseAirDate(s.AirDate),
			MName:        t.Name(),
		})
	}

	tvDetails := TvDetails{
		Id:       tmdbId,
		Name:     res.Name,
		Overview: res.Overview,
		Year:     parseYear(res.FirstAirDate),
		LastEpisodeToAir: TvEpisode{
			Id:            int(res.ID),
			Name:          res.LastEpisodeToAir.Name,
			Overview:      res.LastEpisodeToAir.Overview,
			SeasonNumber:  res.LastEpisodeToAir.SeasonNumber,
			EpisodeNumber: res.LastEpisodeToAir.EpisodeNumber,
			AirDate:       parseAirDate(res.LastEpisodeToAir.AirDate),
			Year:          parseYear(res.LastEpisodeToAir.AirDate),
			MName:         t.Name(),
		},
		NumberOfSeasons:  res.NumberOfSeasons,
		NumberOfEpisodes: res.NumberOfEpisodes,
		Seasons:          seasons,
		MName:            t.Name(),
	}

	if err != nil {
		return nil, err
	}

	return &tvDetails, nil
}

func (t *MetaTmdb) GetTVSeasonDetails(tmdbId int, season int) (*TvSeason, error) {
	res, err := t.c.GetTVSeasonDetails(tmdbId, season, nil)

	if err != nil {
		return nil, err
	}

	tvSeason := TvSeason{
		Id:           int(res.ID),
		MName:        t.Name(),
		Name:         res.Name,
		Overview:     res.Overview,
		SeasonNumber: season,
		EpisodeCount: len(res.Episodes),
		Episodes:     []TvEpisode{},
		AirDate:      parseAirDate(res.AirDate),
	}

	for _, e := range res.Episodes {
		tvSeason.Episodes = append(tvSeason.Episodes, TvEpisode{
			Id:            int(e.ID),
			Name:          e.Name,
			Overview:      e.Overview,
			SeasonNumber:  e.SeasonNumber,
			EpisodeNumber: e.EpisodeNumber,
			Runtime:       e.Runtime,
			AirDate:       parseAirDate(e.AirDate),
		})
	}

	return &tvSeason, err
}

func NewTmdbMeta(c *config.Config) *MetaTmdb {
	m := &MetaTmdb{}
	var err error
	m.c, err = tmdb.Init(c.TmdbApiKey)

	m.c.SetClientAutoRetry()

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5

	m.c.SetClientConfig(*retryClient.StandardClient())

	if err != nil {
		panic(err)
	}

	return m
}
