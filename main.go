package main

import (
	"fmt"
	"log/slog"
	"os"

	"gitlab.com/raffleberry/riptvtime/internal/api"
	"gitlab.com/raffleberry/riptvtime/internal/config"
	"gitlab.com/raffleberry/riptvtime/internal/db"
	"gitlab.com/raffleberry/riptvtime/internal/meta"
	"gitlab.com/raffleberry/riptvtime/internal/services"
	"gitlab.com/raffleberry/riptvtime/internal/utils"
)

func main() {
	logLevel := slog.LevelInfo
	isDev := utils.IsGoRun()
	if isDev {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: isDev,
	}))
	slog.SetDefault(logger)

	var c *config.Config
	var err error

	if isDev {
		c, err = config.LoadFromEnv()
	} else {
		c, err = config.LoadFromFile(os.Args[1])
	}

	if err != nil {
		slog.Error("Failed to load Env", "err", err)
		panic(err)
	}

	addr := fmt.Sprintf("%v:%v", c.Ip, c.Port)
	d := db.NewDbSqlite(c, logger)
	m := meta.NewTmdbMeta(c)

	tvSrv := services.NewTvService(d, m)

	a := api.NewApi(d, m, tvSrv)
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
