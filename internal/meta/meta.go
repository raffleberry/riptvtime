package meta

import (
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("meta not found")
	ErrConfused = errors.New("got more than expected")
)

type TvSearchResult struct {
	Id       int
	Name     string
	Overview string
	Year     int
	Image    string
	Genres   []int64

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
	ShowId        int
	Name          string
	Overview      string
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
	ImgPoster    string

	MName string
}

type TvDetails struct {
	Id               int
	Name             string
	Overview         string
	Year             int
	LastEpisodeToAir TvEpisode
	NextEpisodeToAir TvEpisode
	NumberOfSeasons  int
	NumberOfEpisodes int
	Seasons          []TvSeason
	InProduction     bool
	Tagline          string
	ImgPoster        string
	ImgBackdrop      string
	Genres           []Genre

	MName string
}

type Genre struct {
	Id   int64
	Name string
}

type Meta interface {
	Name() string
	Search(query string, page int) (*TvSearchResults, error)
	GetTvDetails(mId int) (*TvDetails, error)
	GetTVSeasonDetails(mId int, season int) (*TvSeason, error)
	GetTVEpisodeDetails(mId int, season int, episode int) (*TvEpisode, error)
	GetTVFromTvTimeId(tvTimeId int) (*TvDetails, error)
	GetEpisodeFromTvTimeId(tvTimeEId int) (*TvEpisode, error)
	GetGenresTv() ([]Genre, error)
	GetImdbId(mId int) (string, error)
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
