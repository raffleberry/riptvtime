package services

import (
	"archive/zip"
	"encoding/csv"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// tracking-prod-records-v2.csv
type ImportedSeries struct {
	Key string

	Name      string
	TvTimeSId int
	MId       int

	IsStopped bool

	IsFavourite bool

	CreatedAt time.Time
}

type ImportedTrackedEps struct {
	Key string

	SeriesName string
	TvTimeSId  int
	MId        int

	Season  int
	Episode int

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

type Imported struct {
	s *SeriesService
}

func NewImportService(_s *SeriesService) *Imported {
	return &Imported{
		s: _s,
	}
}

func (ipt *Imported) ImportTvTimeSeries(zipPath string) error {

	var allowedFiles = struct {
		SeriesTrackingDataFile string
		FavouriteSeriesFile    string
	}{
		"tracking-prod-records-v2.csv",
		"lists-prod-lists.csv",
	}

	stat, err := os.Stat(zipPath)
	if err != nil {
		return err
	}

	zf, err := os.OpenFile(zipPath, os.O_RDONLY, stat.Mode())

	zr, err := zip.NewReader(zf, stat.Size())
	if err != nil {
		return ErrBadZip
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

	var tsrs []*ImportedSeries
	var teps []*ImportedTrackedEps
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
				return csvErr(fname, err)
			}
			tsrs, teps, err = ipt.ProcessRecsV2(recs)
			if err != nil {
				return csvErr(fname, err)
			}
			processedCnt++
		case allowedFiles.FavouriteSeriesFile:
			recs, err := getRecsFromZf(f)
			if err != nil {
				return csvErr(fname, err)
			}

			favs, err = ipt.ProcessFavs(recs)
			if err != nil {
				return csvErr(fname, err)
			}
			processedCnt++
		}
	}
	ipt.Stub(tsrs, teps, favs)
	allowedFilesCnt := reflect.ValueOf(allowedFiles).NumField()
	if processedCnt != allowedFilesCnt {
		err := errors.Join(ErrBadZip, fmt.Errorf("(Expected %d files but processed %d files) in zip archive", allowedFilesCnt, processedCnt))
		return err
	}

	return nil
}

func (ipt *Imported) Stub(tsrs []*ImportedSeries, teps []*ImportedTrackedEps, favs []int) {
	panic("unimplemented")
}

// lists-prod-lists.csv
func (ipt *Imported) ProcessFavs(lists [][]string) ([]int, error) {

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
func (ipt *Imported) ProcessRecsV2(recsV2 [][]string) ([]*ImportedSeries, []*ImportedTrackedEps, error) {

	hdr := struct {
		SId        string
		SeriesName string

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

	if len(hi) != reflect.ValueOf(hdr).NumField() {
		slog.Error("tracking-prod-records-v2.csv headers", "headers", hi)
		return nil, nil, errors.Join(ErrBadCsv, fmt.Errorf("tracking-prod-records-v2.csv header length(%v)", len(hi)))
	}

	get := func(i int, hname string) string {
		return strings.TrimSpace(recsV2[i][hi[hname]])
	}

	rvSrs := []*ImportedSeries{}
	rvEps := []*ImportedTrackedEps{}

	for i := 1; i < len(recsV2); i++ {
		rec := recsV2[i]
		if len(rec) < len(h) {
			return nil, nil, errors.Join(ErrBadCsv, fmt.Errorf("tracking-prod-records-v2.csv[%v] length(%v)", i, len(rec)))
		}

		key := strings.TrimSpace(rec[hi[hdr.Key]])

		if strings.HasPrefix(key, "watch-episode") || strings.HasPrefix(key, "rewatch-episode") {

			ttsidStr := get(i, hdr.SId)
			ttSid, err := strconv.Atoi(ttsidStr)
			if err != nil {
				slog.Warn("Err parsing csv", "TvTimeSId", ttsidStr, "err", err, "index", i)
			}

			seasonNumberStr := get(i, hdr.SeasonNumber)
			seasonNumber, err := strconv.Atoi(seasonNumberStr)
			if err != nil {
				slog.Warn("Err parsing csv", "Season", seasonNumberStr, "err", err, "index", i)
			}

			sNoStr := get(i, hdr.SNo)
			sNo, err := strconv.Atoi(sNoStr)
			if err != nil {
				slog.Warn("Err parsing csv", "Season", sNoStr, "err", err, "index", i)
			}

			episodeNumberStr := get(i, hdr.EpisodeNumber)
			episodeNumber, err := strconv.Atoi(episodeNumberStr)
			if err != nil {
				slog.Warn("Err parsing csv", "episodeNumber", episodeNumberStr, "err", err, "index", i)
			}

			epNoStr := get(i, hdr.EpNo)
			epNo, err := strconv.Atoi(epNoStr)
			if err != nil {
				slog.Warn("Err parsing csv", "epNo", epNoStr, "err", err, "index", i)
			}

			cAtStr := get(i, hdr.CreatedAt)
			cAt, err := time.Parse(time.DateTime, cAtStr)
			if err != nil {
				slog.Warn("Err parsing csv", "CreatedAt", cAtStr, "err", err, "index", i)
			}

			rvEps = append(rvEps, &ImportedTrackedEps{
				Key:        key,
				SeriesName: get(i, hdr.SeriesName),
				TvTimeSId:  ttSid,
				MId:        0,
				Season:     max(sNo, seasonNumber),
				Episode:    max(epNo, episodeNumber),
				CreatedAt:  cAt,
			})

		} else if strings.HasPrefix(key, "user-series") {

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
				TvTimeSId: ttSid,
				MId:       0,
				IsStopped: isStopped,
				CreatedAt: cAt,
			})
		}
	}

	return rvSrs, rvEps, nil
}
