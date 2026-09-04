package meta_test

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/raffleberry/riptvtime/internal/api"
	"github.com/raffleberry/riptvtime/internal/config"
	"github.com/raffleberry/riptvtime/internal/meta"
)

func Test_ImdbSvc(t *testing.T) {

	startTime := time.Now()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))
	slog.SetDefault(logger)

	mux := http.NewServeMux()
	mux.Handle("GET /", http.FileServer(http.Dir("./testdata")))
	s := api.NewServer("127.0.0.1:0", mux)
	s.Start()
	defer s.Stop()
	slog.Info("imdb test server started", "addr", s.Addr())

	cfg := &config.Config{}
	cfg.ImdbTmpDir = filepath.Join(os.TempDir(), "riptvtime_import__Test_ImdbSvc")
	cfg.ConfigDir = cfg.ImdbTmpDir
	cfg.ImdbDataUrl = s.Addr()
	isvc, err := meta.NewImdbService(logger, cfg)

	gzPath, err := isvc.DownloadRatingsData()
	log.Println(gzPath)
	if err != nil {
		t.Fatal(err)
	}

	err = isvc.DumpImdbRatingsData(gzPath)
	if err != nil {
		t.Fatal(err)
	}

	got := isvc.DbRatingsCnt()
	want := 1712043
	if got != want {
		t.Errorf("want %d, got %v", want, got)
	}

	if !isvc.State.LastUpdated.After(startTime) {
		t.Errorf("want startTime < lastUpdated, got startTime(%v), lastUpdated(%v)", startTime, isvc.State.LastUpdated)
	}

}
