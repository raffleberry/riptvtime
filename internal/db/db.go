package db

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type TvStatus int

var TvStatusVals = []string{"not watching", "watching", "stopped", "completed"}

const (
	TvStatusNotWatching TvStatus = iota
	TvStatusWatching
	TvStatusStopped
	TvStatusCompleted
)

var (
	ErrNotFound = errors.New("Record Not Found")
	ErrBadData  = errors.New("Bad Data")
)

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

type Db interface {
	SeriesTrackedAll() (*[]TvSeries, error)
	SeriesWatchingAll() (*[]TvSeries, error)
	SeriesAdd(t *TvSeries) (int, error)
	SeriesRem(mId int) error

	// 0 if series doesn't exist
	SeriesStatusGet(mId int) (TvStatus, error)
	SeriesStatusUpdate(mId int, newStatus TvStatus) error

	SeriesTrackedEps(mId int) (*[]TvTrackedEps, error)
	SeriesTrackedEpsAdd(ep *TvTrackedEps) (int, error)
	SeriesTrackedEpRemove(mId int, season int, episode int) error

	SeriesSeasonGet(mId int, season int) (*TvSeason, error)
	SeriesSeasonAdd(t *TvSeason) error

	SeriesEpisodeGet(mId int, season int, episode int) (*TvEpisode, error)
}
