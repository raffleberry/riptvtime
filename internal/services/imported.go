package services

import (
	"archive/zip"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/raffleberry/riptvtime/internal/config"
	"github.com/raffleberry/riptvtime/internal/services/state"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
)

var (
	ErrAlreadyProcessing = errors.New("Upload already in progress")
)

// tracking-prod-records-v2.csv
type ImportedSeries struct {
	Key string `gorm:"primaryKey"`

	Name     string
	TvTimeId int
	MId      int

	IsStopped bool

	IsFavourite bool

	Imported      bool `gorm:"default:false"`
	Ignored       bool `gorm:"default:false"`
	Unresolved    bool `gorm:"default:false"`
	UnresolvedMsg string

	UpdatedAt time.Time
	CreatedAt time.Time
}

type ImportedTrackedEps struct {
	Key string `gorm:"primaryKey"`

	SeriesName string
	// series
	TvTimeId int
	// episode
	TvTimeEId int
	MId       int

	Season  int
	Episode int

	Imported      bool `gorm:"default:false"`
	Ignored       bool `gorm:"default:false"`
	Unresolved    bool `gorm:"default:false"`
	UnresolvedMsg string

	UpdatedAt time.Time
	CreatedAt time.Time
}

// tracking-prod-records.csv
type ImportedMovies struct {
	// csv

}

var (
	ErrBadZip = errors.New("Invalid or corrupt zip file")
	ErrBadCsv = errors.New("Invalid or corrupt csv file")
)

type ImportSvc struct {
	ImportTmpDir string
	idb          *gorm.DB
	cfg          *config.Config
}

func NewImportService(logger *slog.Logger, cfg *config.Config) *ImportSvc {

	iTmpDir := cfg.ImportTmpDir
	if iTmpDir == "" {
		tmp := os.TempDir()
		iTmpDir = filepath.Join(tmp, "riptvtime_import_tmp")
	}

	err := os.MkdirAll(iTmpDir, 0755)
	if err != nil {
		panic(err)
	}

	files, err := os.ReadDir(iTmpDir)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		p := filepath.Join(iTmpDir, f.Name())
		if err := os.RemoveAll(p); err != nil {
			slog.Warn("Failed to remove old tmp files: ", "file", p)
		}
	}

	iDbPath := filepath.Join(cfg.ConfigDir, fmt.Sprintf("%s.db", "tvtime_imports"))

	slog.Debug("Initializing import Sqlite Database", "path", iDbPath)

	idb, err := gorm.Open(sqlite.Open(fmt.Sprintf("%v?", iDbPath)), &gorm.Config{
		Logger: glog.NewSlogLogger(logger, glog.Config{
			IgnoreRecordNotFoundError: true,
		}),
	})
	if err != nil {
		panic(err)
	}

	err = idb.AutoMigrate(&ImportedSeries{}, &ImportedTrackedEps{})

	if err != nil {
		panic(err)
	}

	return &ImportSvc{
		ImportTmpDir: iTmpDir,
		idb:          idb,
	}
}

func (ipt *ImportSvc) CleanImportedData() {
	ipt.idb.Unscoped().Where("1=1").Delete(&ImportedSeries{})
	ipt.idb.Unscoped().Where("1=1").Delete(&ImportedTrackedEps{})
}

func (ipt *ImportSvc) isEpisode(key string) bool {
	return strings.HasPrefix(key, "watch-episode") || strings.HasPrefix(key, "rewatch-episode")
}

func (ipt *ImportSvc) isSeries(key string) bool {
	return strings.HasPrefix(key, "user-series")
}

func (ipt *ImportSvc) SetSeriesUnresolved(key string, reason string) error {
	err := ipt.idb.Model(&ImportedSeries{}).Where("key = ?", key).Updates(map[string]any{"unresolved_msg": reason, "unresolved": true}).Error
	return err
}

func (ipt *ImportSvc) SetEpisodeUnresolved(key string, reason string) error {
	err := ipt.idb.Model(&ImportedTrackedEps{}).Where("key = ?", key).Updates(map[string]any{"unresolved_msg": reason, "unresolved": true}).Error
	return err
}

func (ipt *ImportSvc) isDumped(key string) bool {
	var err error

	if ipt.isSeries(key) {
		err = ipt.idb.First(&ImportedSeries{}, "key = ?", key).Error
	} else if ipt.isEpisode(key) {
		err = ipt.idb.First(&ImportedTrackedEps{}, "key = ?", key).Error
	} else {
		return false
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}

	return err == nil
}

func (ipt *ImportSvc) IsImported(key string) bool {
	var err error
	var rv bool
	if ipt.isSeries(key) {
		s := ImportedSeries{}
		err = ipt.idb.First(&s, "key = ?", key).Error
		rv = s.Imported
	} else if ipt.isEpisode(key) {
		e := ImportedTrackedEps{}
		err = ipt.idb.First(&e, "key = ?", key).Error
		rv = e.Imported
	} else {
		return false
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	} else if err != nil {
		log.Panicln(key)
	}

	return rv
}

func (ipt *ImportSvc) SetImported(key string) error {
	if ipt.isSeries(key) {
		return ipt.idb.Model(&ImportedSeries{}).Where("key = ?", key).Update("imported", true).Error
	} else if ipt.isEpisode(key) {
		return ipt.idb.Model(&ImportedTrackedEps{}).Where("key = ?", key).Update("imported", true).Error
	} else {
		return nil
	}
}

func (ipt *ImportSvc) DumpTvTimeGdprData(zipPath string) (int, int, error) {
	var allowedFiles = struct {
		SeriesTrackingDataFile string
		FavouriteSeriesFile    string
	}{
		"tracking-prod-records-v2.csv",
		"lists-prod-lists.csv",
	}

	stat, err := os.Stat(zipPath)
	if err != nil {
		return 0, 0, err
	}

	zf, err := os.OpenFile(zipPath, os.O_RDONLY, stat.Mode())

	zr, err := zip.NewReader(zf, stat.Size())
	if err != nil {
		return 0, 0, ErrBadZip
	}

	defer zf.Close()

	getRecsFromZf := func(f *zip.File) ([][]string, error) {
		fp, err := f.Open()
		if err != nil {
			return nil, err
		}

		csvRdr := csv.NewReader(fp)

		recs, err := csvRdr.ReadAll()
		if err != nil {
			return nil, err
		}
		return recs, nil
	}

	var isrs []*ImportedSeries
	var iteps []*ImportedTrackedEps
	var favs []int

	csvErr := func(fn string, err error) error {
		return errors.Join(ErrBadZip, fmt.Errorf("Error processing CSV target file %s: %v", fn, err))
	}
	processedCnt := 0
	for _, f := range zr.File {
		fname := filepath.Base(f.Name)
		switch fname {
		case allowedFiles.SeriesTrackingDataFile:
			recs, err := getRecsFromZf(f)
			if err != nil {
				return 0, 0, csvErr(fname, err)
			}
			isrs, iteps, err = ipt.ProcessRecsV2(recs)
			if err != nil {
				return 0, 0, csvErr(fname, err)
			}
			processedCnt++
		case allowedFiles.FavouriteSeriesFile:
			recs, err := getRecsFromZf(f)
			if err != nil {
				return 0, 0, csvErr(fname, err)
			}

			favs, err = ipt.ProcessFavs(recs)
			if err != nil {
				return 0, 0, csvErr(fname, err)
			}
			processedCnt++
		}
	}

	allowedFilesCnt := reflect.ValueOf(allowedFiles).NumField()
	if processedCnt != allowedFilesCnt {
		err := errors.Join(ErrBadZip, fmt.Errorf("(Expected %d files but processed %d files) in zip archive", allowedFilesCnt, processedCnt))
		return 0, 0, err
	}

	for _, s := range isrs {
		s.IsFavourite = slices.Contains(favs, s.TvTimeId)
	}

	newSeriesCnt := 0
	newTrackedEpsCnt := 0

	state.Import.UploadingSrsCntTotal = len(isrs)
	state.Import.UploadingEpsCntTotal = len(iteps)

	for _, s := range isrs {
		state.Import.UploadingSrsCnt += 1
		if ipt.isDumped(s.Key) {
			continue
		}
		tx := ipt.idb.Save(s)
		newSeriesCnt += int(tx.RowsAffected)
		if tx.Error != nil {
			return 0, 0, tx.Error
		}
	}

	for _, e := range iteps {
		state.Import.UploadingEpsCnt += 1
		if ipt.isDumped(e.Key) {
			continue
		}
		tx := ipt.idb.Save(e)
		newTrackedEpsCnt += int(tx.RowsAffected)
		if tx.Error != nil {
			return 0, 0, tx.Error
		}
	}

	return newSeriesCnt, newTrackedEpsCnt, nil
}

func (ipt *ImportSvc) GetUnmatched() (*ImportedData, error) {
	var srs []*ImportedSeries
	err := ipt.idb.Where("imported = ?", false).Where("ignored = ?", false).Where("unresolved = ?", false).Find(&srs).Error
	if err != nil {
		return nil, err
	}

	var eps []*ImportedTrackedEps
	err = ipt.idb.Where("imported = ?", false).Where("ignored = ?", false).Where("unresolved = ?", false).Find(&eps).Error
	if err != nil {
		return nil, err
	}
	return &ImportedData{
		Series:   srs,
		Episodes: eps,
	}, nil
}

func (ipt *ImportSvc) GetUnresolved() (*ImportedData, error) {
	var srs []*ImportedSeries
	err := ipt.idb.Where("unresolved = ?", true).Find(&srs).Error
	if err != nil {
		return nil, err
	}

	var eps []*ImportedTrackedEps
	err = ipt.idb.Where("unresolved = ?", true).Find(&eps).Error
	if err != nil {
		return nil, err
	}
	return &ImportedData{
		Series:   srs,
		Episodes: eps,
	}, nil
}

func (ipt *ImportSvc) GetMatched() ([]*ImportedSeries, []*ImportedTrackedEps, error) {
	var rv1 []*ImportedSeries
	err := ipt.idb.Find(&rv1, "m_id != 0").Error
	if err != nil {
		return nil, nil, err
	}

	var rv2 []*ImportedTrackedEps
	err = ipt.idb.Find(&rv2, "m_id != 0").Error
	if err != nil {
		return nil, nil, err
	}
	return rv1, rv2, nil
}

func (ipt *ImportSvc) DeleteSeries(key string) (int, error) {
	tx := ipt.idb.Where("Key = ?", key).Delete(&ImportedSeries{})
	return int(tx.RowsAffected), tx.Error
}

func (ipt *ImportSvc) DeleteTrackedEps(key string) (int, error) {
	tx := ipt.idb.Where("Key = ?", key).Delete(&ImportedTrackedEps{})
	return int(tx.RowsAffected), tx.Error
}

func (ipt *ImportSvc) Match(tvTimeSId, mId int) error {
	tx := ipt.idb.Model(&ImportedSeries{}).Where("tv_time_s_id = ?", tvTimeSId).Update("m_id", mId)
	slog.Debug("ipt matched", "tvTimeSId", tvTimeSId, "mId", mId, "table", "ImportedSeries", "Rows Affected", tx.RowsAffected)
	if tx.Error != nil {
		return tx.Error
	}

	tx = ipt.idb.Model(&ImportedTrackedEps{}).Where("tv_time_s_id = ?", tvTimeSId).Update("m_id", mId)
	slog.Debug("ipt matched", "tvTimeSId", tvTimeSId, "mId", mId, "table", "ImportedTrackedEps", "Rows Affected", tx.RowsAffected)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (ipt *ImportSvc) UpdateMeta(mId int, tvTimeSId int) (int, error) {
	updatedEntriesCnt := 0
	err := ipt.idb.Transaction(func(tx *gorm.DB) error {
		t := tx.Model(&ImportedSeries{}).Where("tv_time_s_id = ?", tvTimeSId).Update("m_id", mId)
		if t.Error != nil {
			return t.Error
		}
		updatedEntriesCnt += int(t.RowsAffected)
		t = tx.Model(&ImportedTrackedEps{}).Where("tv_time_s_id = ?", tvTimeSId).Update("m_id", mId)
		if t.Error != nil {
			return t.Error
		}
		updatedEntriesCnt += int(t.RowsAffected)
		return nil
	})

	return updatedEntriesCnt, err
}

// lists-prod-lists.csv
func (ipt *ImportSvc) ProcessFavs(lists [][]string) ([]int, error) {

	hdr := struct {
		SKey      string
		IsPublic  string
		Type      string
		Objects   string
		CreatedAt string
	}{
		"s_key",
		"is_public",
		"type",
		"objects",
		"created_at",
	}

	if len(lists) < 1 {
		return nil, errors.Join(ErrBadCsv, fmt.Errorf("Empty lists-prod-lists.csv file length(%v)", len(lists)))
	}

	h := lists[0]

	hi := map[string]int{}
	for i, v := range h {
		switch v {
		case hdr.SKey:
			hi[hdr.SKey] = i
		case hdr.IsPublic:
			hi[hdr.IsPublic] = i
		case hdr.Type:
			hi[hdr.Type] = i
		case hdr.Objects:
			hi[hdr.Objects] = i
		case hdr.CreatedAt:
			hi[hdr.CreatedAt] = i
		}
	}

	wantHi := reflect.ValueOf(hdr).NumField()
	if len(hi) != wantHi {
		slog.Error("lists-prod-lists.csv headers", "headers", hi)
		return nil, errors.Join(ErrBadCsv, fmt.Errorf("lists-prod-lists.csv found %v fields but needed %v fields(%v)", len(hi), wantHi))
	}

	favs := []int{}

	for i := 1; i < len(lists); i++ {
		rec := lists[i]
		if len(rec) < len(h) {
			return nil, errors.Join(ErrBadCsv, fmt.Errorf("lists-prod-lists.csv[%v] length(%v)", i, len(rec)))
		}

		sKey := strings.TrimSpace(rec[hi[hdr.SKey]])

		if sKey == "favorite-series" {
			obj := strings.TrimSpace(rec[hi[hdr.Objects]])
			re := regexp.MustCompile(`\s(id\:\d+)`)
			matches := re.FindAllString(obj, -1)
			for _, m := range matches {
				idStr := strings.TrimPrefix(strings.TrimSpace(m), "id:")

				id, err := strconv.Atoi(idStr)
				if err != nil {
					slog.Warn("Bad Match", "m", m, "idStr", idStr, "err", err)
					continue
				}
				favs = append(favs, id)
			}
			break
		}
	}

	return favs, nil
}

// tracking-prod-records-v2.csv
func (ipt *ImportSvc) ProcessRecsV2(recsV2 [][]string) ([]*ImportedSeries, []*ImportedTrackedEps, error) {

	hdr := struct {
		SId        string
		SeriesName string

		EpId      string
		EpisodeId string

		SeasonNumber string
		SNo          string

		EpisodeNumber string
		EpNo          string

		Key string

		IsArchived string
		IsFollowed string

		CreatedAt string
	}{
		"s_id",
		"series_name",

		"ep_id",
		"episode_id",

		"season_number",
		"s_no",

		"episode_number",
		"ep_no",

		"key",

		"is_archived",
		"is_followed",

		"created_at",
	}

	if len(recsV2) < 1 {
		return nil, nil, errors.Join(ErrBadCsv, fmt.Errorf("Empty tracking-prod-records-v2.csv file length(%v)", len(recsV2)))
	}
	h := recsV2[0]

	hi := map[string]int{}
	for i, v := range h {
		switch v {
		case hdr.SId:
			hi[hdr.SId] = i
		case hdr.SeriesName:
			hi[hdr.SeriesName] = i
		case hdr.EpId:
			hi[hdr.EpId] = i
		case hdr.EpisodeId:
			hi[hdr.EpisodeId] = i
		case hdr.SeasonNumber:
			hi[hdr.SeasonNumber] = i
		case hdr.SNo:
			hi[hdr.SNo] = i
		case hdr.EpisodeNumber:
			hi[hdr.EpisodeNumber] = i
		case hdr.EpNo:
			hi[hdr.EpNo] = i
		case hdr.Key:
			hi[hdr.Key] = i
		case hdr.IsArchived:
			hi[hdr.IsArchived] = i
		case hdr.IsFollowed:
			hi[hdr.IsFollowed] = i
		case hdr.CreatedAt:
			hi[hdr.CreatedAt] = i
		}
	}
	hiWant := reflect.ValueOf(hdr).NumField()
	if len(hi) != hiWant {
		slog.Error("tracking-prod-records-v2.csv headers", "headers", hi)
		return nil, nil, errors.Join(ErrBadCsv, fmt.Errorf("tracking-prod-records-v2.csv want-cols#(%v) got-cols#(%v)", hiWant, len(hi)))
	}

	get := func(i int, hname string) string {
		return strings.TrimSpace(recsV2[i][hi[hname]])
	}

	rvSrs := []*ImportedSeries{}
	rvEps := []*ImportedTrackedEps{}

	for i := 1; i < len(recsV2); i++ {
		rec := recsV2[i]
		if len(rec) < len(h) {
			return nil, nil, errors.Join(ErrBadCsv, fmt.Errorf("tracking-prod-records-v2.csv inconsistent row at idx(%v) with length(%v). want(%v)", i, len(rec), hiWant))
		}

		key := strings.TrimSpace(rec[hi[hdr.Key]])

		if ipt.isEpisode(key) {

			ttsidStr := get(i, hdr.SId)
			ttSid, err := strconv.Atoi(ttsidStr)
			if err != nil {
				slog.Warn("parsing csv", "TvTimeSId", ttsidStr, "err", err, "index", i)
			}

			ttEpIdStr := get(i, hdr.EpId)
			ttEpId, errEId1 := strconv.Atoi(ttEpIdStr)

			ttEpisodeIdStr := get(i, hdr.EpisodeId)
			ttEpisodeId, errEId2 := strconv.Atoi(ttEpisodeIdStr)
			if errEId1 != nil && errEId2 != nil {
				slog.Warn("parsing csv", "ttEpId", ttEpIdStr, "errEId1", errEId1, "ttEpisodeId", ttEpisodeIdStr, "errEId2", errEId2, "index", i)
			}

			seasonNumberStr := get(i, hdr.SeasonNumber)
			seasonNumber, errS1 := strconv.Atoi(seasonNumberStr)
			sNoStr := get(i, hdr.SNo)
			sNo, errS2 := strconv.Atoi(sNoStr)
			if errS1 != nil && errS2 != nil {
				slog.Warn("parsing csv", "seasonNumber", seasonNumberStr, "errS1", errS1, "sNo", sNoStr, "errS2", errS2, "index", i)
			}

			episodeNumberStr := get(i, hdr.EpisodeNumber)
			episodeNumber, errE1 := strconv.Atoi(episodeNumberStr)
			epNoStr := get(i, hdr.EpNo)
			epNo, errE2 := strconv.Atoi(epNoStr)
			if errE1 != nil && errE2 != nil {
				slog.Warn("parsing csv", "episodeNumber", episodeNumberStr, "errE1", errE1, "epNo", epNoStr, "errE2", errE2, "index", i)
			}

			cAtStr := get(i, hdr.CreatedAt)
			cAt, err := time.Parse(time.DateTime, cAtStr)
			if err != nil {
				slog.Warn("parsing csv", "CreatedAt", cAtStr, "err", err, "index", i)
			}

			rvEps = append(rvEps, &ImportedTrackedEps{
				Key:        key,
				SeriesName: get(i, hdr.SeriesName),
				TvTimeId:   ttSid,
				TvTimeEId:  max(ttEpId, ttEpisodeId),
				Season:     max(sNo, seasonNumber),
				Episode:    max(epNo, episodeNumber),
				CreatedAt:  cAt,
			})

		} else if ipt.isSeries(key) {

			isFollowedStr := get(i, hdr.IsFollowed)
			isFollowed, err := strconv.ParseBool(isFollowedStr)
			if err != nil {
				slog.Warn("Err parsing csv", "IsFollowed", isFollowedStr, "err", err, "index", i)
			}

			if !isFollowed {
				continue
			}

			isArchivedStr := get(i, hdr.IsArchived)
			isArchived, err := strconv.ParseBool(isArchivedStr)
			if err != nil {
				slog.Warn("Err parsing csv", "IsArchived", isArchivedStr, "err", err, "index", i)
			}

			isStopped := isArchived

			ttsidStr := get(i, hdr.SId)
			ttSid, err := strconv.Atoi(ttsidStr)
			if err != nil {
				slog.Warn("Err parsing csv", "TvTimeSId", ttsidStr, "err", err, "index", i)
			}

			cAtStr := get(i, hdr.CreatedAt)
			cAt, err := time.Parse(time.DateTime, cAtStr)
			if err != nil {
				slog.Warn("Err parsing csv", "CreatedAt", cAtStr, "err", err, "index", i)
			}

			rvSrs = append(rvSrs, &ImportedSeries{
				Key:       key,
				Name:      get(i, hdr.SeriesName),
				TvTimeId:  ttSid,
				MId:       0,
				IsStopped: isStopped,
				CreatedAt: cAt,
			})
		}
	}

	return rvSrs, rvEps, nil
}
