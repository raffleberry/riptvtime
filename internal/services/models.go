package services

import (
	"time"

	"github.com/raffleberry/riptvtime/internal/db"
	"github.com/raffleberry/riptvtime/internal/meta"
)

type SeriesFeedItem struct {
	db.TvSeries
	EpisodesTotal   int
	EpisodesAired   int
	EpisodesWatched int
	UpNext          *SeriesEpisode
	RecentlyAired   bool
	Image           string

	LastEpisodeAirDate   time.Time
	LastEpisodeWatchDate time.Time
	ShowAddDate          time.Time

	ImdbRating *meta.ImdbRating
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
	S         int
	E         int
	Cnt       int
	CreatedAt time.Time
}

type SeriesFullItem struct {
	*meta.TvDetails
	ImdbRating    *meta.ImdbRating
	EpisodesAired int
	// including special episodes(specials)
	EpsWatched []SeriesEpisode
}

type SeriesTracked struct {
	db.TvSeries
	InProduction bool
}

type ImportedData struct {
	Series   []*ImportedSeries
	Episodes []*ImportedTrackedEps
}

type SeriesFavs struct {
	db.TvSeriesFav
	ImgPoster string
}

type UpcomingItem struct {
	SeriesName string
	Year       int
	ImgPoster  string
	Episode    *db.TvEpisode
	ImdbRating *meta.ImdbRating
}

type Genre struct {
	meta.Genre
	Cnt int
}

// Discover list names served by SeriesService.Discover
const (
	DiscoverRecommended = "recommended"
	DiscoverAiringToday = "airing_today"
	DiscoverOnTheAir    = "on_the_air"
	DiscoverTopRated    = "top_rated"
	DiscoverPopular     = "popular"
)

var DiscoverLists = []string{
	DiscoverRecommended,
	DiscoverAiringToday,
	DiscoverOnTheAir,
	DiscoverTopRated,
	DiscoverPopular,
}

var DiscoverSortOptions = []string{
	"popularity.desc",
	"popularity.asc",
	"first_air_date.desc",
	"first_air_date.asc",
	"vote_average.desc",
	"vote_average.asc",
	"vote_count.desc",
	"vote_count.asc",
	"name.asc",
	"name.desc",
	"original_name.asc",
	"original_name.desc",
}
