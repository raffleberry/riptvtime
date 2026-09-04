package meta

import (
	"compress/gzip"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/raffleberry/riptvtime/internal/config"
)

var (
	ErrAlreadyRefreshing = errors.New("Refresh already in progress")
	ErrUnAvailable       = errors.New("Meta Unavailable")
)

const (
	ImdbRatingsTbl = `
	CREATE TABLE IF NOT EXISTS imdb_ratings(
		id TEXT PRIMARY KEY,
		rating INTEGER,
		votes INTEGER
	);
	`
	ImdbStateTbl = `
	CREATE TABLE IF NOT EXISTS imdb_state(
		ky TEXT PRIMARY KEY,
		last_updated DATETIME,
		refresh_every INTEGER,
		not_found_cnt INTEGER
	);
	`
)

var tables = []string{ImdbRatingsTbl, ImdbStateTbl}

type ImdbState struct {
	Ky          string
	LastUpdated time.Time
	// days
	RefreshEvery int
	NotFoundCnt  int
}

type ImdbMeta struct {
	refreshing bool
	rmu        sync.Mutex
	mu         sync.Mutex
	tmpDir     string
	db         *sql.DB
	cfg        *config.Config
	dataPk     string
	dataUrl    string

	// do not modify
	State *ImdbState
}

// TODO: Convert this to return err,
// this should be non-essential that wont break the app on startup
func NewImdbService(logger *slog.Logger, cfg *config.Config) (*ImdbMeta, error) {
	iTmpDir := cfg.ImdbTmpDir
	if iTmpDir == "" {
		return nil, fmt.Errorf("No tmp dir provided for imdb downloads")
	}
	err := os.MkdirAll(iTmpDir, 0755)
	if err != nil {
		return nil, err
	}

	imdbDataUrl := cfg.ImdbDataUrl
	if imdbDataUrl == "" {
		imdbDataUrl = "https://datasets.imdbws.com"
	}

	files, err := os.ReadDir(iTmpDir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		p := filepath.Join(iTmpDir, f.Name())
		if err := os.RemoveAll(p); err != nil {
			slog.Warn("Failed to remove old tmp files: ", "file", p)
		}
	}

	err = os.MkdirAll(cfg.ConfigDir, 0755)
	if err != nil {
		slog.Error("Failed to ensure config dir exists", "err", err)
		return nil, err
	}

	// No difference between WAL and OFF
	// iDbPath := filepath.Join(cfg.ConfigDir, fmt.Sprintf("%s.db?_pragma=journal_mode(WAL)&_pragma=synchronous(OFF)", "imdb"))
	iDbPath := filepath.Join(cfg.ConfigDir, fmt.Sprintf("%s.db?", "imdb"))

	slog.Debug("Initializing imdb Sqlite Database", "path", iDbPath)

	db, err := sql.Open("sqlite", iDbPath)

	if err != nil {
		slog.Error("Failed to connect imdb local db")
		return nil, err
	}

	for i := range tables {
		_, err := db.Exec(tables[i])
		if err != nil {
			slog.Error("Failed to create imdb table", "err", err, "query", tables[i])
			return nil, err
		}
	}

	dataPk := "imdb"
	data := ImdbState{
		Ky:           dataPk,
		LastUpdated:  time.Time{},
		RefreshEvery: 7,
		NotFoundCnt:  0,
	}

	row := db.QueryRow(`SELECT last_updated, refresh_every, not_found_cnt FROM imdb_state WHERE ky = ?`, dataPk)
	err = row.Scan(&data.LastUpdated, &data.RefreshEvery, &data.NotFoundCnt)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = db.Exec(`INSERT INTO imdb_state(ky, last_updated, refresh_every, not_found_cnt) VALUES (?, ?, ?, ?)`, dataPk, data.LastUpdated, data.RefreshEvery, data.NotFoundCnt)
		if err != nil {
			slog.Error("Failed to insert imdb data state", "err", err)
			return nil, err
		}
	} else if err != nil {
		slog.Error("Failed to get info about imdb data state")
		return nil, err
	}

	return &ImdbMeta{
		tmpDir:  iTmpDir,
		db:      db,
		cfg:     cfg,
		dataPk:  dataPk,
		dataUrl: imdbDataUrl,
		State:   &data,
	}, nil
}

func (svc *ImdbMeta) GetRating(id string) (ImdbRating, error) {
	rv := ImdbRating{}

	row := svc.db.QueryRow(`SELECT id, rating, votes FROM imdb_ratings WHERE id = ?`, id)
	err := row.Scan(&rv.Id, &rv.Rating, &rv.Votes)

	if errors.Is(err, sql.ErrNoRows) {
		svc.IncrNotFoundCnt()
		if svc.State.NotFoundCnt != 0 {
			days := time.Since(svc.State.LastUpdated).Hours() / 24
			if !svc.refreshing && days > float64(svc.State.RefreshEvery) {
				go svc.Refresh()
			}
			if err != nil {
				return rv, err
			}
			svc.State.NotFoundCnt = 0
		}
		return rv, ErrNotFound
	} else if err != nil {
		return rv, err
	}

	return rv, nil
}

func (svc *ImdbMeta) Refresh() error {
	svc.rmu.Lock()
	svc.refreshing = true
	defer func() {
		svc.refreshing = false
		svc.rmu.Unlock()
	}()
	path, err := svc.DownloadRatingsData()
	if err != nil {
		slog.Error("REFRESH(imdbSvc): Failed to download imdb ratings data", "err", err)
		return err
	}
	err = svc.DumpImdbRatingsData(path)
	if err != nil {
		slog.Error("REFRESH(imdbSvc): Failed to dump imdb ratings data", "err", err)
		return err
	}
	return svc.ResetNotFoundCnt()
}

func (svc *ImdbMeta) IncrNotFoundCnt() {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	_, err := svc.db.Exec(`UPDATE imdb_state SET not_found_cnt=?`, 1+svc.State.NotFoundCnt)
	if err != nil {
		slog.Error("Failed to increment Not Found Cnt", "err", err)
		return
	}
	svc.State.NotFoundCnt++
}

func (svc *ImdbMeta) ResetNotFoundCnt() error {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	_, err := svc.db.Exec(`UPDATE imdb_state SET not_found_cnt=?`, 0)
	if err != nil {
		slog.Error("Failed to reset NotFoundCnt", "err", err)
		return err
	}
	svc.State.NotFoundCnt = 0
	return nil
}

func (svc *ImdbMeta) DbRatingsCnt() int {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	cnt := 0
	row := svc.db.QueryRow(`SELECT count(*) FROM imdb_ratings`)
	err := row.Scan(&cnt)
	if err != nil {
		slog.Error("Failed to get imdb ratings count", "err", err)
		return -1
	}

	return cnt
}

func (svc *ImdbMeta) DumpImdbRatingsData(gzPath string) error {

	stat, err := os.Stat(gzPath)
	if err != nil {
		return err
	}

	gzf, err := os.OpenFile(gzPath, os.O_RDONLY, stat.Mode())

	gzReader, err := gzip.NewReader(gzf)
	if err != nil {
		return ErrBadGZip
	}

	defer gzf.Close()

	tsvReader := csv.NewReader(gzReader)
	defer gzReader.Close()

	tsvReader.Comma = '\t'
	if _, err := tsvReader.Read(); err != nil {
		slog.Error("Failed to read header", "err", err)
		return err
	}

	count := 0

	tx, err := svc.db.Begin()
	if err != nil {
		slog.Error("Failed to begin transaction", "err", err)
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`DROP TABLE IF EXISTS imdb_ratings`)
	if err != nil {
		slog.Error("Failed to drop `imdb_ratings` table")
		return fmt.Errorf("Failed to drop `imdb_ratings` table: %v", err)
	}

	_, err = tx.Exec(ImdbRatingsTbl)
	if err != nil {
		slog.Error("Failed to create table", "err", err)
		return fmt.Errorf("Failed to create table: %v", err)
	}

	slog.Info("Imdb inserting ratings table...")
	stmt, err := tx.Prepare(`INSERT INTO imdb_ratings(id, rating, votes) VALUES(?, ?, ?)`)
	if err != nil {
		slog.Error("Failed to prepare statement", "err", err)
		return err
	}
	defer stmt.Close()

	for {
		record, err := tsvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Warning: skipping corrupt line: %v", err)
			continue
		}

		id := record[0]
		rating, _ := (strconv.ParseFloat(record[1], 64))
		votes, _ := strconv.Atoi(record[2])

		_, err = stmt.Exec(id, rating, votes)
		if err != nil {
			slog.Error("Failed to insert record", "err", err)
			return err
		}

		count++

		if count%100000 == 0 {
			slog.Debug("Inserted", "count", count)
		}
	}

	slog.Debug("Successfully inserted imdb ratings!", "count", count)

	slog.Debug("Updating LastInserted")

	svc.mu.Lock()
	defer svc.mu.Unlock()

	lastUpdated := time.Now()
	svc.State.LastUpdated = lastUpdated

	_, err = tx.Exec(`UPDATE imdb_state SET last_updated=?`, time.Now())
	if err != nil {
		slog.Error("Failed to set last_updated", "err", err)
		return err
	}

	err = tx.Commit()

	if err != nil {
		slog.Error("Failed to commit transaction", "err", err)
		return err
	}

	svc.State.LastUpdated = lastUpdated

	return nil
}

func (svc *ImdbMeta) DownloadRatingsData() (string, error) {
	destTsv := filepath.Join(svc.tmpDir, fmt.Sprintf("title.ratings.%d.tsv.gz", time.Now().Unix()))
	tsvGz, err := os.Create(destTsv)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer tsvGz.Close()

	url := fmt.Sprintf("%s/title.ratings.tsv.gz", svc.dataUrl)
	slog.Info("imdb ratings file", "url", url)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status code: %s", resp.Status)
	}

	_, err = io.Copy(tsvGz, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save file contents: %w", err)
	}

	return destTsv, nil
}
