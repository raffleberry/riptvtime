package db

import (
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"

	"gitlab.com/raffleberry/riptvtime/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	glog "gorm.io/gorm/logger"
)

type DbSqlite struct {
	orm *gorm.DB
	err error
}

func NewDbSqlite(c *config.Config, logger *slog.Logger) *DbSqlite {
	db := &DbSqlite{}
	sqliteDbPath := filepath.Join(c.ConfigDir, "riptvtime.db")

	slog.Debug("Initializing Sqlite Database", "path", sqliteDbPath)

	db.orm, db.err = gorm.Open(sqlite.Open(fmt.Sprintf("%v?", sqliteDbPath)), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger: glog.NewSlogLogger(logger, glog.Config{
			IgnoreRecordNotFoundError: true,
		}),
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

func (db *DbSqlite) SeriesWatchingAll() (*[]TvSeries, error) {
	var series []TvSeries
	err := db.orm.Where("tracking_status = ?", TvStatusWatching).Find(&series).Error
	return &series, err
}

func (db *DbSqlite) SeriesTrackedEps(tmdbId int) (*[]TvTrackedEps, error) {
	var trackedEps []TvTrackedEps
	err := db.orm.Where("series_m_id = ?", tmdbId).Find(&trackedEps).Error
	return &trackedEps, err
}

func (db *DbSqlite) SeriesAdd(t *TvSeries) (int, error) {
	err := db.orm.Create(t).Error

	return int(t.ID), err
}

func (db *DbSqlite) SeriesRem(id int) error {
	err := db.orm.Delete(&TvSeries{}, "m_id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
	}

	return err
}

func (db *DbSqlite) SeriesSeasonAdd(t *TvSeason) (int, error) {
	err := db.orm.Create(t).Error

	return int(t.ID), err
}

func (db *DbSqlite) SeriesStatusGet(mId int) (TvStatus, error) {
	var res TvSeries
	err := db.orm.Model(&TvSeries{}).Select("tracking_status").Where("m_id = ?", mId).Limit(1).First(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return TvStatusNotWatching, nil
		}
	}
	return res.TrackingStatus, err
}

func (db *DbSqlite) SeriesEpisodeGet(mId int, season int, episode int) (*TvEpisode, error) {
	var ep TvEpisode
	err := db.orm.Take(&ep, "series_m_id = ? AND season = ? AND episode = ?", mId, season, episode).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &ep, err
}

// (insertId, err)
func (db *DbSqlite) SeriesTrackedEpsAdd(ep *TvTrackedEps) (int, error) {
	err := db.orm.Create(ep).Error
	return int(ep.ID), err
}

// (rowsAffected, err)
func (db *DbSqlite) SeriesTrackedEpRemove(mId int, season int, episode int) (int, error) {
	tx := db.orm.Where("series_m_id = ? AND season = ? AND episode = ?", mId, season, episode).
		Limit(1).
		Delete(&TvTrackedEps{})

	if tx.Error != nil {
		return 0, tx.Error
	}

	if tx.RowsAffected == 0 {
		return 0, ErrNotFound
	}

	return int(tx.RowsAffected), nil
}

func (db *DbSqlite) SeriesStatusUpdate(mId int, newStatus TvStatus) error {
	err := db.orm.Model(&TvSeries{}).Where("m_id = ?", mId).Update("tracking_status", newStatus).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
	}

	return err
}

type CacheSqlite struct {
	orm *gorm.DB
	err error
}

func NewCacheSqlite(cfg *config.Config, logger *slog.Logger) *CacheSqlite {
	c := CacheSqlite{}

	cachePath := filepath.Join(cfg.ConfigDir, "cache.db")

	c.orm, c.err = gorm.Open(sqlite.Open(fmt.Sprintf("%v?", cachePath)), &gorm.Config{
		Logger: glog.NewSlogLogger(logger, glog.Config{}),
	})

	if c.err != nil {
		panic("failed to connect cache database")
	}

	err := c.orm.AutoMigrate(&Cached{})

	if err != nil {
		slog.Error("Failed to migrate cache", "err", err)
		panic("Failed to migrate database")
	}

	return &c
}

func (c *CacheSqlite) Set(data *Cached) error {
	return c.orm.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"json_data"}),
	}).Create(data).Error

}

func (c *CacheSqlite) Get(key string) (*Cached, error) {
	var cached Cached
	err := c.orm.Where("key = ?", key).First(&cached).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &cached, err

}
