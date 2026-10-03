package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/raffleberry/riptvtime/internal/api"
	"github.com/raffleberry/riptvtime/internal/config"
	"github.com/raffleberry/riptvtime/internal/db"
	"github.com/raffleberry/riptvtime/internal/meta"
	"github.com/raffleberry/riptvtime/internal/services"
	"github.com/raffleberry/riptvtime/internal/setup"
	"github.com/raffleberry/riptvtime/internal/utils"
)

func main() {

	prefix := flag.String("prefix", "", "server prefix")
	logLvl := flag.Int("log-level", 8, "value of log level = {-4, 0, 4, 8 <- default} | where -4 = debug, 0 = info, 4 = warn, 8 = error")
	conf := flag.String("conf", "", "pass config file path")
	flag.Parse()

	logLevel := slog.Level(*logLvl)

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
		if len(*conf) == 0 {
			cfg, err = setup.GetConfigFromUser()
			if cfg == nil {
				if err != nil {
					slog.Error("Failed to get config", "err", err)
				}
				return
			}
		} else {
			cfg, err = config.LoadFromFile(*conf)
		}
	}

	if err != nil {
		slog.Error("Failed to load config", "err", err)
		panic(err)
	}

	addr := fmt.Sprintf("%v:%v", cfg.Ip, cfg.Port)
	d := db.NewDbSqlite(cfg, logger)
	im, err := meta.NewImdbService(logger, cfg)
	m := meta.NewTmdbMeta(cfg, im, logger)
	iptSrv := services.NewImportService(logger, cfg)
	if err != nil {
		slog.Error("Failed to start imdb service. Disabling Imdb")
		cfg.EnableImdb = false
	}
	tvSrv := services.NewTvService(cfg, d, m, iptSrv)

	a := api.NewApi(*prefix, d, m, tvSrv, cfg)
	s := api.NewServer(addr, a.Router)

	fmt.Printf("Starting server...\n")

	if err := s.Start(); err != nil {
		panic(err)
	}

	url := fmt.Sprintf("http://%s/%s", addr, strings.TrimPrefix(*prefix, "/"))
	if !setup.BrowserOpened {
		utils.OpenBrowser(url)
		setup.BrowserOpened = true
	}
	fmt.Printf("Address - %s ...\n", url)
	s.WaitSIGINT()

	fmt.Printf("Stopping http://%s/ ...\n", addr)
	s.Stop()

}
