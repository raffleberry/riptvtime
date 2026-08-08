package db

import (
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/raffleberry/riptvtime/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	glog "gorm.io/gorm/logger"
)

func GetInProdExpireTime() time.Time {
	Hours := rand.Int63n(60-36+1) + 36
	return time.Now().Add(time.Duration(Hours) * time.Hour)
}

func GetNotInProdExpireTime() time.Time {
	days14_21 := rand.Int63n(21-14+1) + 14
	return time.Now().Add(time.Duration(days14_21) * time.Hour * 24)
}

type Cached struct {
	Key       string `gorm:"primaryKey"`
	JsonData  string
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiredAt time.Time
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
		DoUpdates: clause.AssignmentColumns([]string{"json_data", "updated_at", "expired_at"}),
	}).Create(data).Error
}

func (c *CacheSqlite) Get(key string) (*Cached, error) {
	var cached Cached
	err := c.orm.Where("key = ?", key).First(&cached).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %w: key=%v", err, ErrNotFound, key)
		}
		return nil, err
	}
	return &cached, err

}
