package db

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"gitlab.com/raffleberry/riptvtime/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Db struct {
	orm *gorm.DB
	err error
}

func NewDb(c *config.Config) *Db {
	db := &Db{}
	sqliteDbPath := filepath.Join(c.ConfigDir, "riptvtime.db")

	slog.Debug("Initializing Sqlite Database", "path", sqliteDbPath)

	db.orm, db.err = gorm.Open(sqlite.Open(fmt.Sprintf("%v?", sqliteDbPath)), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if db.err != nil {
		panic("failed to connect database")
	}

	err := db.orm.AutoMigrate(&TvSeries{}, &TvTracking{}, &TvSeason{}, &TvEpisode{})

	if err != nil {
		slog.Error("Failed to migrate", "err", err)
		panic("Failed to migrate database")
	}

	return db
}

// func NewTrackTv(tmdbId int) (uint, error) {
// 	td, err := TmdbClient.GetTVDetails(tmdbId, nil)

// 	if err != nil {
// 		return 0, err
// 	}

// 	ts := TvSeries{
// 		TmdbId:         td.ID,
// 		Name:           td.Name,
// 		Overview:       td.Overview,
// 		Genres:         genreToStr(td.Genres),
// 		Year:           ParseYear(td.FirstAirDate),
// 		FirstAirDate:   ParseAirDate(td.FirstAirDate),
// 		TrackingStatus: TvStatusWatching,
// 	}

// 	res := Db.Create(&ts)

// 	if res.Error != nil {
// 		return 0, res.Error
// 	}

// 	return ts.ID, nil
// }
