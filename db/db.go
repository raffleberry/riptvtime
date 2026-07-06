package db

import (
	"log/slog"
	"path/filepath"

	"gitlab.com/raffleberry/riptvtime/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Conn *gorm.DB
var err error
var SQLITE_DB_PATH string

func init() {
	initTmdb()
	initDb()
}

func initDb() {
	SQLITE_DB_PATH = filepath.Join(utils.APP_CONF_DIR, "riptvtime.db")

	slog.Debug("Initializing Sqlite Database", "path", SQLITE_DB_PATH)

	Conn, err = gorm.Open(sqlite.Open(SQLITE_DB_PATH), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	Conn.AutoMigrate(&TvSeries{})
	Conn.AutoMigrate(&TvTracking{})

}
