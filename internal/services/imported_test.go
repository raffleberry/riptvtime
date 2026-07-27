package services_test

import (
	"encoding/csv"
	"os"
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
	isrs, ieps, err := ipt.ProcessRecsV2(recsV2)

	if err != nil {
		t.Error(err)
	}

	type x struct{}

	requiredSrs := map[string]struct{}{
		utils.Jn("The Blacklist", 266189):             x{},
		utils.Jn("Barry", 333072):                     x{},
		utils.Jn("Wrecked (2016)", 310555):            x{},
		utils.Jn("Luther", 159591):                    x{},
		utils.Jn("You", 336924):                       x{},
		utils.Jn("Avatar: The Last Airbender", 74852): x{},
	}

	for _, isr := range isrs {
		if _, ok := requiredSrs[utils.Jn(isr.Name, isr.TvTimeSId)]; !ok {
			t.Errorf("Unexpected series: %v", isr.Name)
		}
	}

	requiredEpsCnt := map[string]int{
		"You":                        30,
		"The Blacklist":              10,
		"Barry":                      17,
		"Wrecked (2016)":             15,
		"Luther":                     7,
		"Avatar: The Last Airbender": 41,
	}

	epsCnt := map[string]int{}

	for _, ie := range ieps {
		epsCnt[ie.SeriesName]++
	}

	for k, v := range requiredEpsCnt {
		if v != epsCnt[k] {
			t.Errorf("Unexpected ep count for %s: %d != %d", k, v, epsCnt[k])
		}
	}
}
