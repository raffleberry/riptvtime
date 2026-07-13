package services

import (
	"gitlab.com/raffleberry/riptvtime/internal/db"
	"gitlab.com/raffleberry/riptvtime/internal/meta"
)

type SeriesFeedItem struct {
	db.TvSeries
	EpisodesTotal   int
	EpisodesAired   int
	EpisodesWatched int
	UpNextS         int
	UpNextE         int
	RecentlyAired   bool
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
