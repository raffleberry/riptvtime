package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/raffleberry/riptvtime/internal/api"
	"github.com/raffleberry/riptvtime/internal/config"
	"github.com/raffleberry/riptvtime/internal/db"
	"github.com/raffleberry/riptvtime/internal/meta"
	"github.com/raffleberry/riptvtime/internal/services"
	"github.com/raffleberry/riptvtime/internal/setup"
	"github.com/raffleberry/riptvtime/internal/utils"
)

func main() {
	logLevel := slog.LevelInfo
	isDev := utils.IsGoRun()
	if isDev {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: isDev,
	}))
	slog.SetDefault(logger)

	var cfg *config.Config
	var err error

	if isDev {
		cfg, err = config.LoadFromEnv()
	} else {
		if len(os.Args) < 2 {
			cfg, err = setup.GetConfigFromUser()
			if cfg == nil {
				if err != nil {
					slog.Error("Failed to get config", "err", err)
				}
				return
			}
		} else {
			cfg, err = config.LoadFromFile(os.Args[1])
		}
	}

	if err != nil {
		slog.Error("Failed to load config", "err", err)
		panic(err)
	}

	addr := fmt.Sprintf("%v:%v", cfg.Ip, cfg.Port)
	d := db.NewDbSqlite(cfg, logger)
	m := meta.NewTmdbMeta(cfg)
	iptSrv := services.NewImportService(logger, cfg)
	tvSrv := services.NewTvService(d, m, iptSrv)

	a := api.NewApi(d, m, tvSrv, cfg)
	s := api.NewServer(addr, a.Router)

	fmt.Printf("Starting server...\n")

	if err := s.Start(); err != nil {
		panic(err)
	}

	fmt.Printf("Address - http://%s/ ...\n", addr)
	s.WaitSIGINT()

	fmt.Printf("Stopping http://%s/ ...\n", addr)
	s.Stop()

}
