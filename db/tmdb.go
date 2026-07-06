package db

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/hashicorp/go-retryablehttp"
)

var TmdbClient *tmdb.Client

func initTmdb() {
	var err error
	TmdbClient, err = tmdb.Init(os.Getenv("TMDB_API_KEY"))

	TmdbClient.SetClientAutoRetry()

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5

	TmdbClient.SetClientConfig(*retryClient.StandardClient())

	if err != nil {
		panic(err)
	}

}

func GenreToStr(genres []tmdb.Genre) string {
	var r strings.Builder
	for i, g := range genres {
		if i != 0 {
			r.WriteString(",")
		}
		r.WriteString(g.Name)
	}
	return r.String()
}

func ParseYear(date string) int {
	if len(date) >= 4 {
		var err error
		year, err := strconv.Atoi(date[:4])
		if err != nil {
			slog.Warn("Error parsing year", "date", date, "err", err)
			year = 0
		}
		return year
	}
	slog.Warn("Failed to parse year", "date", date)
	return 0
}

func ParseAirDate(date string) time.Time {
	if len(date) == 10 {
		firstAirDate, err := time.Parse(time.DateOnly, date)
		if err != nil {
			slog.Warn("Error parsing AirDate", "date", date, "err", err)
		}
		return firstAirDate
	}
	slog.Warn("Failed to parse date", "date", date)
	return time.Time{}
}
