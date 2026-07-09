package db

import (
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"

	"gitlab.com/raffleberry/riptvtime/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DbSqlite struct {
	orm *gorm.DB
	err error
}

func NewDbSqlite(c *config.Config) *DbSqlite {
	db := &DbSqlite{}
	sqliteDbPath := filepath.Join(c.ConfigDir, "riptvtime.db")

	slog.Debug("Initializing Sqlite Database", "path", sqliteDbPath)

	db.orm, db.err = gorm.Open(sqlite.Open(fmt.Sprintf("%v?", sqliteDbPath)), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if db.err != nil {
		panic("failed to connect database")
	}

	err := db.orm.AutoMigrate(&TvSeries{}, &TvTrackedEps{}, &TvSeason{}, &TvEpisode{})

	if err != nil {
		slog.Error("Failed to migrate", "err", err)
		panic("Failed to migrate database")
	}

	return db
}

func (db *DbSqlite) SeriesWatching() (*[]TvSeries, error) {
	var series []TvSeries
	err := db.orm.Find(&series).Where("tracking_status = ?", TvStatusWatching).Error
	return &series, err
}

func (db *DbSqlite) SeriesTrackedEps(tmdbId int) (*[]TvTrackedEps, error) {
	var trackedEps []TvTrackedEps
	err := db.orm.Find(&trackedEps).Where("series_tmdb_id = ?", tmdbId).Error
	return &trackedEps, err
}

func (db *DbSqlite) SeriesAdd(t *TvSeries) (int, error) {
	err := db.orm.Create(t).Error

	return int(t.ID), err
}

func (db *DbSqlite) SeriesSeasonAdd(t *TvSeason) (int, error) {
	err := db.orm.Create(t).Error

	return int(t.ID), err
}

func (db *DbSqlite) SeriesIsAdded(mId int) (bool, error) {
	var added bool
	err := db.orm.Model(&TvSeries{}).Select("count(*) > 0").Where("m_id = ?", mId).Limit(1).Find(&added).Error
	return added, err
}

func (db *DbSqlite) SeriesEpisodeExists(mId int, season int, episode int) (bool, error) {
	var exists bool
	err := db.orm.Model(&TvEpisode{}).Select("count(*) > 0").Where("series_m_id = ? AND season = ? AND episode = ?", mId, season, episode).Limit(1).Find(&exists).Error
	return exists, err
}

func (db *DbSqlite) SeriesEpisodeGet(mId int, season int, episode int) (*TvEpisode, error) {
	var ep TvEpisode
	err := db.orm.Take(&ep, "series_m_id = ? AND season = ? AND episode = ?", mId, season, episode).Error
	return &ep, err
}

// (insertId, err)
func (db *DbSqlite) SeriesTrackedEpsAdd(ep *TvTrackedEps) (int, error) {
	err := db.orm.Create(ep).Error
	return int(ep.ID), err
}

// (rowsAffected, err)
func (db *DbSqlite) SeriesTrackedEpRemove(mId int, season int, episode int) (int, error) {
	ep := TvTrackedEps{}
	err := db.orm.First(&ep, "series_m_id = ? AND season = ? AND episode = ?", mId, season, episode).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrNotFound
		}
		return 0, err
	}

	tx := db.orm.Delete(&ep)
	if tx.Error != nil {
		return 0, tx.Error
	}

	return int(tx.RowsAffected), nil
}
