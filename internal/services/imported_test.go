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

	got, err := ipt.ProcessFavs(recsV2)
	if err != nil {
		t.Error(err)
	}
	wantIds := []int{75897, 79335, 82459, 95491,
		176941, 281662, 293088, 94571, 362472, 79349,
		244061, 74852, 425999, 332331, 366211, 452595}
	if len(wantIds) != len(got.Series) {
		t.Errorf("Unexpected  series count: want(%d) != got(%d)", len(wantIds), len(got.Series))
	}
	for _, s := range got.Series {
		if !slices.Contains(wantIds, s.TvTimeId) {
			t.Errorf("Want Series MId: %d IN got: [%v]", s.TvTimeId, got.Series)
		}
	}

	wantUuids := []string{
		"1be8d227-5d39-4561-8dfa-7520b8c51d0f",
		"c969e184-ae5d-4766-8f94-561aa5f942ef",
		"8c4fffb7-5729-43e1-b10b-886fe62b75b5",
		"e16e4a64-2ec9-4363-b468-11d40fad789e",
		"01d4581e-5764-4597-aad9-46fadc47ae62",
		"bd2e6693-3d05-44eb-b8a8-68ff96b636cd",
		"d1bcb2ec-7a87-4686-9af1-44f4510fe80f",
		"b6e6e8d7-15f5-419e-8f5d-8c9d33c8320f",
		"73e49a28-9aab-4cbe-ab04-3d9cbb778fee",
		"b7c2098d-e182-4ffb-942c-8a81b1c2bf0b",
	}
	if len(wantUuids) != len(got.Movies) {
		t.Errorf("Unexpected  Movies count: want(%d) != got(%d)", len(wantUuids), len(got.Movies))
	}
	for _, s := range got.Movies {
		if !slices.Contains(wantUuids, s.Uuid) {
			t.Errorf("Want Movies Uuid: %s IN got: [%v]", s.Uuid, got.Movies)
		}
	}

}

func Test_DumpTvTimeGdprData(t *testing.T) {
	const zfp = "testdata/tvtime-gdpr-data.zip"
	cfg := &config.Config{}
	cfg.ImportTmpDir = filepath.Join(os.TempDir(), "riptvtime_import__Test_DumpTvTimeGdprData")
	cfg.ConfigDir = cfg.ImportTmpDir
	ipt := services.NewImportService(nil, cfg)
	res, err := ipt.DumpTvTimeGdprData(zfp)

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
	wantEpsCnt := map[string]int{
		"You":                        30,
		"The Blacklist":              10,
		"Barry":                      17,
		"Wrecked (2016)":             15,
		"Luther":                     7,
		"Avatar: The Last Airbender": 41,
	}
	wantFavSrs, wantFavMovs := 16, 10
	epsTotCnt := 0
	for _, v := range wantEpsCnt {
		epsTotCnt += v
	}
	if res.SeriesCnt != len(wantSrs) {
		t.Errorf("wantSrsCnt(%v) != gotSrsCnt(%v)", len(wantSrs), res.SeriesCnt)
	}

	if res.EpisodesCnt != epsTotCnt {
		t.Errorf("wantEpsCnt(%v) != gotEpsCnt(%v)", epsTotCnt, res.EpisodesCnt)
	}
	if res.FavSeriesCnt != wantFavSrs {
		t.Errorf("wantFavSrs(%v) != gotFavSrs(%v)", wantFavSrs, res.FavSeriesCnt)
	}
	if res.FavMoviesCnt != wantFavMovs {
		t.Errorf("wantFavMovs(%v) != gotFavMovs(%v)", wantFavMovs, res.FavMoviesCnt)
	}

	// Dump again, but rows affected should be 0
	res, err = ipt.DumpTvTimeGdprData(zfp)
	if err != nil {
		t.Error(err)
	}
	if res.SeriesCnt != 0 {
		t.Errorf("wantSrsCnt(%v) != gotSrsCnt(%v)", 0, res.SeriesCnt)
	}
	if res.EpisodesCnt != 0 {
		t.Errorf("wantEpsCnt(%v) != gotEpsCnt(%v)", 0, res.EpisodesCnt)
	}
	if res.FavSeriesCnt != 0 {
		t.Errorf("wantFavSrs(%v) != gotFavSrs(%v)", 0, res.FavSeriesCnt)
	}
	if res.FavMoviesCnt != 0 {
		t.Errorf("wantFavMovs(%v) != gotFavMovs(%v)", 0, res.FavMoviesCnt)
	}

}
