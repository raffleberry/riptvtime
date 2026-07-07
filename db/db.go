package db

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"gitlab.com/raffleberry/riptvtime/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB
var err error
var SQLITE_DB_PATH string

func init() {
	initTmdb()
	initDb()
}

func initDb() {
	SQLITE_DB_PATH = filepath.Join(utils.APP_CONF_DIR, "riptvtime.db")

	slog.Debug("Initializing Sqlite Database", "path", SQLITE_DB_PATH)

	Db, err = gorm.Open(sqlite.Open(fmt.Sprintf("%v?", SQLITE_DB_PATH)), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("failed to connect database")
	}

	err := Db.AutoMigrate(&TvSeries{}, &TvTracking{}, &TvSeason{}, &TvEpisode{})

	if err != nil {
		slog.Error("Failed to migrate", "err", err)
		panic("Failed to migrate database")
	}

}

func NewTrackTv(tmdbId int) (uint, error) {
	td, err := TmdbClient.GetTVDetails(tmdbId, nil)

	if err != nil {
		return 0, err
	}

	ts := TvSeries{
		TmdbId:         td.ID,
		Name:           td.Name,
		Overview:       td.Overview,
		Genres:         GenreToStr(td.Genres),
		Year:           ParseYear(td.FirstAirDate),
		FirstAirDate:   ParseAirDate(td.FirstAirDate),
		TrackingStatus: TvStatusWatching,
	}

	res := Db.Create(&ts)

	if res.Error != nil {
		return 0, res.Error
	}

	return ts.ID, nil
}
