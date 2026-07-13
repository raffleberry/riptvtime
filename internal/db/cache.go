package db

import (
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"gitlab.com/raffleberry/riptvtime/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	glog "gorm.io/gorm/logger"
)

type Cached struct {
	Key       string `gorm:"primaryKey"`
	JsonData  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Cache interface {
	Get(key string) (*Cached, error)
	Set(data *Cached) error
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
