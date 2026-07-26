package services

import (
	"time"

	"gitlab.com/raffleberry/riptvtime/internal/db"
	"gitlab.com/raffleberry/riptvtime/internal/meta"
)

type SeriesFeedItem struct {
	db.TvSeries
	EpisodesTotal      int
	EpisodesAired      int
	EpisodesWatched    int
	UpNextS            int
	UpNextE            int
	RecentlyAired      bool
	LastEpisodeAirDate time.Time
}

type SeriesSearchItem struct {
	meta.TvSearchResult
	Status db.TvStatus
}
type SeriesSearchResult struct {
	Page         int
	Results      []SeriesSearchItem
	TotalPages   int
	TotalResults int
}

type SeriesEpisode struct {
	S int
	E int
}

type SeriesFullItem struct {
	*meta.TvDetails
	EpisodesAired int
	EpsWatched    []SeriesEpisode
}

type SeriesTracked struct {
	db.TvSeries
	InProduction bool
}
