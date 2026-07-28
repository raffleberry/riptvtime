package services_test

import (
	"encoding/csv"
	"os"
	"slices"
	"testing"

	"github.com/raffleberry/riptvtime/internal/services"
	"github.com/raffleberry/riptvtime/internal/utils"
)

func TestProcessRecsV2(t *testing.T) {
	const filePath = "testdata/tracking-prod-records-v2.csv"
	csvFp, err := os.Open(filePath)
	if err != nil {
		t.Error(err)
	}
	csvRdr := csv.NewReader(csvFp)

	recsV2, err := csvRdr.ReadAll()
	if err != nil {
		t.Error(err)
	}

	ipt := services.NewImportService(nil)
	gotSrs, gotEps, err := ipt.ProcessRecsV2(recsV2)

	if err != nil {
		t.Error(err)
	}

	type x struct{}

	wantSrs := map[string]struct{}{
		utils.Jn("The Blacklist", 266189):             x{},
		utils.Jn("Barry", 333072):                     x{},
		utils.Jn("Wrecked (2016)", 310555):            x{},
		utils.Jn("Luther", 159591):                    x{},
		utils.Jn("You", 336924):                       x{},
		utils.Jn("Avatar: The Last Airbender", 74852): x{},
	}

	for _, sr := range gotSrs {
		if _, ok := wantSrs[utils.Jn(sr.Name, sr.TvTimeSId)]; !ok {
			t.Errorf("Unexpected series: %v", sr.Name)
		}
	}

	wantEpsCnt := map[string]int{
		"You":                        30,
		"The Blacklist":              10,
		"Barry":                      17,
		"Wrecked (2016)":             15,
		"Luther":                     7,
		"Avatar: The Last Airbender": 41,
	}

	gotEpsCnt := map[string]int{}

	for _, ie := range gotEps {
		gotEpsCnt[ie.SeriesName]++
	}

	for k, v := range wantEpsCnt {
		if v != gotEpsCnt[k] {
			t.Errorf("Unexpected ep count for %s: %d != %d", k, v, gotEpsCnt[k])
		}
	}
}

func TestProcessFavsLists(t *testing.T) {
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

	ipt := services.NewImportService(nil)
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
