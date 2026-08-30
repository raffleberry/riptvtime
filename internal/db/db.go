package db

import (
	"errors"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

type TvStatus int

var TvStatusVals = []string{"not watching", "watching", "stopped", "completed", "up to date"}

const (
	// when not found in db
	TvStatusNotWatching TvStatus = iota

	TvStatusWatching
	TvStatusStopped

	// Derived Status {if in_production = false}
	TvStatusCompleted
	// Derived Status {if in_production = true}
	TvStatusUpToDate
)

const (
	SourceImport string = "import"
	SourceUI     string = "ui"
)

var (
	ErrNotFound = errors.New("Record Not Found")
	ErrBadData  = errors.New("Bad Data")
)

func GetInProdExpireTime() time.Time {
	Hours := rand.Int63n(60-36+1) + 36
	return time.Now().Add(time.Duration(Hours) * time.Hour)
}

func GetNotInProdExpireTime() time.Time {
	days14_21 := rand.Int63n(21-14+1) + 14
	return time.Now().Add(time.Duration(days14_21) * time.Hour * 24)
}

func (s TvStatus) String() string {
	if s < TvStatusWatching || int(s) >= len(TvStatusVals) {
		return "unknown"
	}
	return TvStatusVals[s]
}

func (s TvStatus) IsValid() bool {
	if s < TvStatusNotWatching || int(s) >= len(TvStatusVals) {
		return false
	}
	return true
}

type TvSeries struct {
	gorm.Model
	MName          string
	MId            int64
	Name           string
	Overview       string
	Year           int
	TrackingStatus TvStatus `gorm:"default:1"`
	RuntimeApprox  int
	InProduction   bool `gorm:"default:true"`

	Source    string
	SourceKey string
}

type TvTrackedEps struct {
	gorm.Model
	MName      string
	EpisodeMId int64
	SeriesMId  int64
	Name       string
	Overview   string
	Season     int
	Episode    int
	Runtime    int

	Source    string
	SourceKey string
}

type TvEpisode struct {
	gorm.Model
	MName     string `gorm:"uniqueIndex:idx_episode_mname_mid"`
	MId       int64  `gorm:"uniqueIndex:idx_episode_mname_mid"`
	SeriesMId int64
	Name      string
	Overview  string
	Season    int
	Episode   int
	Runtime   int
	AirDate   time.Time
}

type TvSeason struct {
	gorm.Model
	MName     string `gorm:"uniqueIndex:idx_season_mname_mid"`
	MId       int64  `gorm:"uniqueIndex:idx_season_mname_mid"`
	SeriesMId int64
	AirDate   time.Time
	Season    int
	Name      string
	Overview  string
	Episodes  []TvEpisode `gorm:"foreignKey:SeriesMId,Season;references:SeriesMId,Season"`
}

type TvSeriesDetails struct {
	gorm.Model
	MName    string
	MId      int64
	Name     string
	Overview string
	AirDate  time.Time
}

type Genres struct {
	gorm.Model
	Name string
	MId  int64
}

type Cached struct {
	What      string `gorm:"primaryKey;index:idx_what_key,unique"`
	Key       string `gorm:"primaryKey;index:idx_what_key,unique"`
	JsonData  string
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiredAt time.Time
}

type Stats struct {
	TotalHours    float64
	TotalEpisodes int
	TotalShows    int
}

type StatsShow struct {
	MId  int
	Name string
	Year int

	Image string

	LastActivity string
}

type Db interface {
	SeriesTrackedAll() (*[]TvSeries, error)
	SeriesFeed() ([]TvSeries, error)
	SeriesWatchingInProdAll() ([]TvSeries, error)
	SeriesAdd(t *TvSeries) (int, error)
	SeriesRem(mId int) error

	// 0 if series doesn't exist
	SeriesStatusGet(mId int) (TvStatus, error)
	SeriesStatusUpdate(mId int, newStatus TvStatus) error

	SeriesTrackedEps(mId int) ([]TvTrackedEps, error)
	SeriesTrackedEpsAdd(ep *TvTrackedEps) (int, error)
	SeriesTrackedEpRemove(mId int, season int, episode int) error

	SeriesSeasonGet(mId int, season int) (*TvSeason, error)
	SeriesSeasonAdd(t *TvSeason) error

	SeriesUpdateInProd(id int, inProd bool) error

	SeriesEpisodeGet(mId int, season int, episode int) (*TvEpisode, error)

	ImportedSeriesCheck(key string) (bool, error)
	ImportedTrackedEpsCheck(key string) (bool, error)

	SeriesStats() (*Stats, error)
	SeriesStatsMyShows(limit int) ([]StatsShow, error)

	CacheSet(data *Cached) error
	CacheGet(what, key string) (*Cached, error)
}
