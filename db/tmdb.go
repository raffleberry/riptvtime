package db

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/hashicorp/go-retryablehttp"
)

var TmdbClient *tmdb.Client

func initTmdb() {
	var err error
	TmdbClient, err = tmdb.Init(os.Getenv("TMDB_API_KEY"))

	TmdbClient.SetClientAutoRetry()

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5

	TmdbClient.SetClientConfig(*retryClient.StandardClient())

	if err != nil {
		panic(err)
	}

}

func GenreToStr(genres []tmdb.Genre) string {
	var r strings.Builder
	for i, g := range genres {
		if i != 0 {
			r.WriteString(",")
		}
		r.WriteString(g.Name)
	}
	return r.String()
}

func ParseYear(date string) int {
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

func ParseAirDate(date string) time.Time {
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

func CacheSeason(tmdbId int64, season int) error {
	sd, err := TmdbClient.GetTVSeasonDetails(int(tmdbId), season, nil)
	slog.Debug("Caching Tv Season", "name", sd.Name, "tmdbId", tmdbId, "Season", season, "Episodes", len(sd.Episodes))

	if err != nil {
		return err
	}

	episodes := []TvEpisode{}
	for _, e := range sd.Episodes {
		episodes = append(episodes, TvEpisode{
			TmdbID:     e.ID,
			SeriesTmdb: e.ShowID,
			Name:       e.Name,
			Overview:   e.Overview,
			Season:     e.SeasonNumber,
			Episode:    e.EpisodeNumber,
			Runtime:    e.Runtime,
			AirDate:    ParseAirDate(e.AirDate),
		})
	}

	dbSd := TvSeason{
		TmdbID:     sd.ID,
		AirDate:    ParseAirDate(sd.AirDate),
		SeriesTmdb: tmdbId,
		Season:     season,
		Name:       sd.Name,
		Overview:   sd.Overview,
		Episodes:   episodes,
	}

	res := Db.Create(&dbSd)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func GetEpisodeDetails(seriesTmdbId int64, sNo int, epNo int) (*TvEpisode, error) {
	var exists bool

	err := Db.Model(&TvEpisode{}).
		Select("count(*) > 0").
		Where("series_tmdb = ? AND season = ? AND episode = ?", seriesTmdbId, sNo, epNo).
		Limit(1).
		Find(&exists).Error

	if err != nil {
		slog.Error("Failed to check if episode exists", "err", err)
		return nil, err
	}

	if !exists {
		slog.Debug("Season details not found in local db, fetching from tmdb", "series_tmdb", seriesTmdbId, "season", sNo)
		err := CacheSeason(seriesTmdbId, sNo)
		if err != nil {
			return nil, err
		}
	}

	var ep TvEpisode
	err = Db.Take(&ep, "series_tmdb = ? AND season = ? AND episode = ?", seriesTmdbId, sNo, epNo).Error
	if err != nil {
		slog.Error("Failed to get episode details", "err", err)
		return nil, err
	}

	return &ep, nil
}
