package meta

import (
	"time"
)

type TvSearchResult struct {
	Id       int
	Name     string
	Overview string
	Year     int

	MName string
}

type TvSearchResults struct {
	Page         int
	Results      []TvSearchResult
	TotalPages   int
	TotalResults int
}
type TvEpisode struct {
	Id            int
	Name          string
	Overview      string
	Year          int
	SeasonNumber  int
	EpisodeNumber int
	AirDate       time.Time
	Runtime       int

	MName string
}

type TvSeason struct {
	Id           int
	Name         string
	Overview     string
	SeasonNumber int
	EpisodeCount int
	Episodes     []TvEpisode
	AirDate      time.Time

	MName string
}

type TvDetails struct {
	Id               int
	Name             string
	Overview         string
	Year             int
	LastEpisodeToAir TvEpisode
	NumberOfSeasons  int
	NumberOfEpisodes int
	Seasons          []TvSeason

	MName string
}

type Meta interface {
	Name() string
	Search(query string, page int) (*TvSearchResults, error)
	GetTvDetails(mId int) (*TvDetails, error)
	GetTVSeasonDetails(mId int, season int) (*TvSeason, error)
}

// func GetEpisodeDetails(seriesTmdbId int64, sNo int, epNo int) (*TvEpisode, error) {
// 	var exists bool

// 	err := Db.Model(&TvEpisode{}).
// 		Select("count(*) > 0").
// 		Where("series_tmdb = ? AND season = ? AND episode = ?", seriesTmdbId, sNo, epNo).
// 		Limit(1).
// 		Find(&exists).Error

// 	if err != nil {
// 		slog.Error("Failed to check if episode exists", "err", err)
// 		return nil, err
// 	}

// 	if !exists {
// 		slog.Debug("Season details not found in local db, fetching from tmdb", "series_tmdb", seriesTmdbId, "season", sNo)
// 		err := CacheSeason(seriesTmdbId, sNo)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	var ep TvEpisode
// 	err = Db.Take(&ep, "series_tmdb = ? AND season = ? AND episode = ?", seriesTmdbId, sNo, epNo).Error
// 	if err != nil {
// 		slog.Error("Failed to get episode details", "err", err)
// 		return nil, err
// 	}

// 	return &ep, nil
// }
