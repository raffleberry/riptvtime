package db

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type TvStatus int

var (
	ErrNotFound = errors.New("Record Not Found")
)

const (
	TvStatusWatching TvStatus = iota
	TvStatusStopped
	TvStatusCompleted
)

func (s TvStatus) String() string {
	strings := [...]string{"not watching", "watching", "stopped", "completed"}
	if s < TvStatusWatching || int(s) >= len(strings) {
		return "unknown"
	}
	return strings[s]
}

type TvSeries struct {
	gorm.Model
	MName          string
	MId            int64
	Name           string
	Overview       string
	Year           int
	TrackingStatus TvStatus
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
	MName     string
	MId       int64
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
	MName     string
	MId       int64
	SeriesMId int64
	AirDate   time.Time
	Season    int
	Name      string
	Overview  string
	Episodes  []TvEpisode `gorm:"foreignkey:SeriesMId;references:SeriesMId"`
}

type Db interface {
	SeriesWatching() (*[]TvSeries, error)
	SeriesAdd(t *TvSeries) (int, error)
	SeriesIsAdded(mId int) (bool, error)

	SeriesTrackedEps(mId int) (*[]TvTrackedEps, error)
	SeriesTrackedEpsAdd(ep *TvTrackedEps) (int, error)
	SeriesTrackedEpRemove(mId int, season int, episode int) (int, error)

	SeriesSeasonAdd(t *TvSeason) (int, error)

	SeriesEpisodeGet(mId int, season int, episode int) (*TvEpisode, error)
	SeriesEpisodeExists(mId int, season int, episode int) (bool, error)
}
