package db

import (
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"github.com/raffleberry/riptvtime/internal/config"
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
		panic(db.err)
	}

	err := db.orm.AutoMigrate(&TvSeries{}, &TvTrackedEps{}, &TvSeason{}, &TvEpisode{})

	if err != nil {
		slog.Error("Failed to migrate", "err", err)
		panic("Failed to migrate database")
	}

	return db
}

func (db *DbSqlite) SeriesTrackedAll() (*[]TvSeries, error) {
	var series []TvSeries
	err := db.orm.Find(&series).Error
	return &series, err
}

func (db *DbSqlite) SeriesFeed() ([]TvSeries, error) {
	var series []TvSeries
	err := db.orm.Where("tracking_status = ?", TvStatusWatching).Find(&series).Error
	return series, err
}

func (db *DbSqlite) SeriesWatchingInProdAll() ([]TvSeries, error) {
	var series []TvSeries
	err := db.orm.Where("tracking_status = ?", TvStatusWatching).Where("in_production = ?", true).Find(&series).Error
	return series, err
}

func (db *DbSqlite) SeriesUpdateInProd(id int, inProd bool) error {
	return db.orm.Model(&TvSeries{}).Where("id = ?", id).Update("in_production", inProd).Error
}

func (db *DbSqlite) SeriesTrackedEps(tmdbId int) ([]TvTrackedEps, error) {
	var trackedEps []TvTrackedEps
	err := db.orm.Where("series_m_id = ?", tmdbId).Order("created_at DESC").Find(&trackedEps).Error
	return trackedEps, err
}

func (db *DbSqlite) SeriesAdd(t *TvSeries) (int, error) {
	err := db.orm.Create(t).Error

	return int(t.ID), err
}

func (db *DbSqlite) SeriesRem(id int) error {
	err := db.orm.Delete(&TvSeries{}, "m_id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: Series with MId %d Not Found", ErrNotFound, id)
		}
	}

	return err
}

// upserts
func (db *DbSqlite) SeriesSeasonAdd(t *TvSeason) error {
	return db.orm.Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "m_name"},
				{Name: "m_id"},
			},
			UpdateAll: true,
		}).Create(t).Error
		if err != nil {
			return err
		}

		for _, episode := range t.Episodes {
			// episode.SeriesMId = t.SeriesMId

			err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{
					{Name: "m_name"},
					{Name: "m_id"},
				},
				UpdateAll: true,
			}).Create(&episode).Error
			if err != nil {
				return err
			}
		}

		return nil
	})

}

func (db *DbSqlite) SeriesSeasonGet(mId, season int) (*TvSeason, error) {
	sn := &TvSeason{}
	err := db.orm.Preload("Episodes").Take(&sn, "series_m_id = ? AND season = ?", mId, season).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w: Series with mId=`%v`, season=`%v` not found", err, ErrNotFound, mId, season)
	}
	return sn, err
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
		return nil, fmt.Errorf("%w: %w: Series with mId=`%v`, season=`%v`, episode=`%v` not found", err, ErrNotFound, mId, season, episode)
	}
	return &ep, err
}

// (insertId, err)
func (db *DbSqlite) SeriesTrackedEpsAdd(ep *TvTrackedEps) (int, error) {
	err := db.orm.Create(ep).Error
	return int(ep.ID), err
}

// (rowsAffected, err)
func (db *DbSqlite) SeriesTrackedEpRemove(mId int, season int, episode int) error {
	var ep TvTrackedEps

	err := db.orm.Last(&ep, "series_m_id = ? AND season = ? AND episode = ?", mId, season, episode).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	} else if err != nil {
		return err
	}

	return db.orm.Delete(&ep).Error
}

func (db *DbSqlite) SeriesStatusUpdate(mId int, newStatus TvStatus) error {

	if newStatus != TvStatusStopped && newStatus != TvStatusWatching {
		return nil
	}

	err := db.orm.Model(&TvSeries{}).Where("m_id = ?", mId).Update("tracking_status", newStatus).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: Series with MId %d Not Found", ErrNotFound, mId)
		}
	}

	return err
}

func (db *DbSqlite) ImportedSeriesCheck(key string) (bool, error) {
	var cnt int64
	err := db.orm.Model(&TvSeries{}).Where("source_key = ?", key).Where("source = ?", SourceImport).Count(&cnt).Error
	return cnt > 0, err
}

func (db *DbSqlite) ImportedTrackedEpsCheck(key string) (bool, error) {
	var cnt int64
	err := db.orm.Model(&TvTrackedEps{}).Where("source_key = ?", key).Where("source = ?", SourceImport).Count(&cnt).Error
	return cnt > 0, err
}

func (db *DbSqlite) SeriesStatsTotal() (int, error) {
	var res struct {
		Total int
	}
	err := db.orm.Model(&TvTrackedEps{}).Select("sum(runtime) as total").First(&res).Error
	return res.Total, err
}
