package services_test

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/raffleberry/riptvtime/internal/config"
	"github.com/raffleberry/riptvtime/internal/services"
	"github.com/raffleberry/riptvtime/internal/utils"
)

func Test_ProcessFavs(t *testing.T) {
	const filePath = "testdata/lists-prod-lists.csv"
	csvFp, err := os.Open(filePath)
	if err != nil {
		t.Error(err)
	}
	csvRdr := csv.NewReader(csvFp)

	recsV2, err := csvRdr.ReadAll()
	if err != nil {
		t.Error(err)
	}

	cfg := &config.Config{}
	cfg.ImportTmpDir = filepath.Join(os.TempDir(), "riptvtime_import__Test_ProcessFavs")
	cfg.ConfigDir = cfg.ImportTmpDir
	ipt := services.NewImportService(nil, cfg)

	gotIds, err := ipt.ProcessFavs(recsV2)

	if err != nil {
		t.Error(err)
	}

	type x struct{}

	wantIds := []int{75897, 79335, 82459, 95491,
		176941, 281662, 293088, 94571, 362472, 79349,
		244061, 74852, 425999, 332331, 366211, 452595}

	if len(wantIds) != len(gotIds) {
		t.Errorf("Unexpected series count: want(%d) != got(%d)", len(wantIds), gotIds[0])
	}

	for _, id := range gotIds {
		if !slices.Contains(wantIds, id) {
			t.Errorf("Unexpected series: %v", id)
		}
	}
}

func Test_DumpTvTimeGdprData(t *testing.T) {
	const zfp = "testdata/tvtime-gdpr-data.zip"
	cfg := &config.Config{}
	cfg.ImportTmpDir = filepath.Join(os.TempDir(), "riptvtime_import__Test_DumpTvTimeGdprData")
	cfg.ConfigDir = cfg.ImportTmpDir
	ipt := services.NewImportService(nil, cfg)
	gotSrsCnt, gotEpsCnt, err := ipt.DumpTvTimeGdprData(zfp)

	if err != nil {
		t.Error(err)
		return
	}

	wantSrs := map[string]struct{}{
		utils.Jn("The Blacklist", 266189):             {},
		utils.Jn("Barry", 333072):                     {},
		utils.Jn("Wrecked (2016)", 310555):            {},
		utils.Jn("Luther", 159591):                    {},
		utils.Jn("You", 336924):                       {},
		utils.Jn("Avatar: The Last Airbender", 74852): {},
	}
	if gotSrsCnt != len(wantSrs) {
		t.Errorf("wantSrsCnt(%v) != gotSrsCnt(%v)", len(wantSrs), gotSrsCnt)
	}

	wantEpsCnt := map[string]int{
		"You":                        30,
		"The Blacklist":              10,
		"Barry":                      17,
		"Wrecked (2016)":             15,
		"Luther":                     7,
		"Avatar: The Last Airbender": 41,
	}
	epsTotCnt := 0
	for _, v := range wantEpsCnt {
		epsTotCnt += v
	}
	if epsTotCnt != gotEpsCnt {
		t.Errorf("wantEpsCnt(%v) != gotEpsCnt(%v)", epsTotCnt, gotEpsCnt)

	}

	// Dump again, but rows affected should be 0
	gotSrsCnt, gotEpsCnt, err = ipt.DumpTvTimeGdprData(zfp)
	if err != nil {
		t.Error(err)
	}
	if gotSrsCnt != 0 {
		t.Errorf("wantSrsCnt(%v) != gotSrsCnt(%v)", 0, gotSrsCnt)
	}
	if gotEpsCnt != 0 {
		t.Errorf("wantEpsCnt(%v) != gotEpsCnt(%v)", 0, gotEpsCnt)
	}

	// for _, sr := range gotSrs {
	// 	if _, ok := wantSrs[utils.Jn(sr.Name, sr.TvTimeSId)]; !ok {
	// 		t.Errorf("Unexpected series: %v", sr.Name)
	// 	}
	// 	if sr.TvTimeSId == 74852 && !sr.IsFavourite || sr.TvTimeSId != 74852 && sr.IsFavourite {
	// 		t.Errorf("Unexpected favourite series, only 74852 should be favourite: %v", sr.Name)
	// 	}
	// }

	// gotEpsCnt := map[string]int{}
	// for _, ie := range gotEps {
	// 	gotEpsCnt[ie.SeriesName]++
	// }

	// for k, v := range wantEpsCnt {
	// 	if v != gotEpsCnt[k] {
	// 		t.Errorf("Unexpected ep count for %s: %d != %d", k, v, gotEpsCnt[k])
	// 	}
	// }

}
