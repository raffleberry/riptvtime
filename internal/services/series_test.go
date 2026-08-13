package services_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/raffleberry/riptvtime/internal/db"
	"github.com/raffleberry/riptvtime/internal/meta"
	"github.com/raffleberry/riptvtime/internal/services"
	"gorm.io/gorm"
)

func mockDb() db.Db {
	type MockDb struct {
		db.Db
	}
	return MockDb{}
}

func TestSeriesService_MakeFeedList(t *testing.T) {

	atime, _ := time.Parse(time.DateOnly, "2000-01-01")
	btime, _ := time.Parse(time.DateOnly, "2001-01-01")
	// ctime, _ := time.Parse(time.DateOnly, "2002-01-01")

	getMockEpsWatched := func(seasonEpCnt, tillSeason, tillEpisode int, t time.Time) []services.SeriesEpisode {
		var rv []services.SeriesEpisode
		for i := 1; i <= tillSeason; i++ {
			till := seasonEpCnt
			if i == tillSeason {
				till = tillEpisode
			}
			for j := 1; j <= till; j++ {
				rv = append(rv, services.SeriesEpisode{S: i, E: j, Cnt: 1, CreatedAt: t})
				t = t.Add(time.Hour)
			}
		}
		return rv
	}

	series1 := []db.TvSeries{
		db.TvSeries{
			MId:            1,
			MName:          "TvMetaService",
			Name:           "Test1",
			Overview:       "Overview1",
			TrackingStatus: db.TvStatusWatching,
			Year:           2001,
			Model:          gorm.Model{CreatedAt: atime},
		},
		db.TvSeries{
			MId:            2,
			MName:          "TvMetaService",
			Name:           "Test2",
			Overview:       "Overview2",
			TrackingStatus: db.TvStatusWatching,
			Year:           2002,
			Model:          gorm.Model{CreatedAt: btime},
		},
	}

	freshSeriesData1 := []*services.SeriesFullItem{
		// 1 (in prod, watching)
		&services.SeriesFullItem{
			&meta.TvDetails{
				Id:       1,
				Name:     "UpdatedTest1",
				MName:    "TvMetaService",
				Overview: "UpdatedOverview1",
				Year:     2011,
				LastEpisodeToAir: meta.TvEpisode{
					SeasonNumber:  2,
					EpisodeNumber: 7,
					AirDate:       atime.Add(21 * 24 * time.Hour),
				},
				Seasons: []meta.TvSeason{
					meta.TvSeason{
						Id:           1,
						Name:         "Season 1",
						Overview:     "Overview",
						SeasonNumber: 1,
						EpisodeCount: 13,
					},
					meta.TvSeason{
						Id:           2,
						Name:         "Season 2",
						Overview:     "Overview",
						SeasonNumber: 2,
						EpisodeCount: 13,
					},
				},
				InProduction:     true,
				NumberOfEpisodes: 26,
			},
			20,
			getMockEpsWatched(13, 2, 1, atime.Add((1+21)*24*time.Hour)),
		},
		// 2 (in prod, stopped)
		&services.SeriesFullItem{
			&meta.TvDetails{
				Id:       2,
				Name:     "UpdatedTest2",
				Overview: "UpdatedOverview2",
				Year:     2012,
				LastEpisodeToAir: meta.TvEpisode{
					SeasonNumber:  1,
					EpisodeNumber: 5,
					AirDate:       btime.Add(6 * 24 * time.Hour),
				},
				Seasons: []meta.TvSeason{
					meta.TvSeason{
						Id:           1,
						Name:         "Season 1",
						Overview:     "Overview",
						SeasonNumber: 1,
						EpisodeCount: 13,
					},
				},
				InProduction:     true,
				NumberOfEpisodes: 13,
			},
			5,
			getMockEpsWatched(13, 1, 2, btime.Add((1+6)*24*time.Hour)),
		},
	}

	want1 := []*services.SeriesFeedItem{

		&services.SeriesFeedItem{
			db.TvSeries{
				MId:            2,
				MName:          "TvMetaService",
				Name:           "UpdatedTest2",
				Overview:       "UpdatedOverview2",
				TrackingStatus: db.TvStatusWatching,
				Year:           2012,
			},
			13,
			5,
			2,
			nil,
			false,
			"",
			btime.Add(6 * 24 * time.Hour),
			btime.Add(168 * time.Hour),
			btime,
		},
		&services.SeriesFeedItem{
			db.TvSeries{
				MId:            1,
				MName:          "TvMetaService",
				Name:           "UpdatedTest1",
				Overview:       "UpdatedOverview1",
				TrackingStatus: db.TvStatusWatching,
				Year:           2011,
			},
			26,
			20,
			14,
			nil,
			false,
			"",
			atime.Add(21 * 24 * time.Hour),
			atime.Add(528 * time.Hour),
			atime,
		},
	}

	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		db   db.Db
		meta meta.Meta
		ipt  *services.ImportSvc
		// Named input parameters for target function.
		series          []db.TvSeries
		freshSeriesData []*services.SeriesFullItem
		want            []*services.SeriesFeedItem
	}{
		{
			"Validations - only TvStatusWatching(both in-prod & out-of-prod), TvShowEps Cnt, WatchCount",
			nil, nil, nil,
			series1,
			freshSeriesData1,
			want1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := services.NewTvService(tt.db, tt.meta, tt.ipt)
			got := srv.MakeFeedList(tt.series, tt.freshSeriesData)
			if len(tt.want) != len(got) {
				t.Errorf("MakeFeedList() = got len(%v) != want len(%v)", len(got), len(tt.want))
			}
			for i, g := range got {

				if g.LastEpisodeAirDate != tt.want[i].LastEpisodeAirDate {
					t.Errorf("MakeFeedList() = got LastEpisodeAirDate(%v) != want LastEpisodeAirDate(%v)", g.LastEpisodeAirDate, tt.want[i].LastEpisodeAirDate)
				}
				if g.LastEpisodeWatchDate != tt.want[i].LastEpisodeWatchDate {
					t.Errorf("MakeFeedList() = got LastEpisodeWatchDate(%v) != want LastEpisodeWatchDate(%v)", g.LastEpisodeWatchDate, tt.want[i].LastEpisodeWatchDate)
				}
				if g.ShowAddDate != tt.want[i].ShowAddDate {
					t.Errorf("MakeFeedList() = got ShowAddDate(%v) != want ShowAddDate(%v)", g.ShowAddDate, tt.want[i].ShowAddDate)
				}

				if g.MId != tt.want[i].MId {
					fmt.Printf("%v %v %v %v\n", g.MId, tt.want[i].MId, g.Name, tt.want[i].Name)
					t.Errorf("MakeFeedList() = got MId(%v) != want MId(%v)", g.MId, tt.want[i].MId)
				}
				if g.EpisodesAired != tt.want[i].EpisodesAired {
					t.Errorf("MakeFeedList() = got EpisodesAired(%v) != want EpisodesAired(%v)", g.EpisodesAired, tt.want[i].EpisodesAired)
				}
				if g.EpisodesWatched != tt.want[i].EpisodesWatched {
					t.Errorf("MakeFeedList() = got EpisodesWatched(%v) != want EpisodesWatched(%v)", g.EpisodesWatched, tt.want[i].EpisodesWatched)
				}
				if g.EpisodesTotal != tt.want[i].EpisodesTotal {
					t.Errorf("MakeFeedList() = got EpisodesTotal(%v) != want EpisodesTotal(%v)", g.EpisodesTotal, tt.want[i].EpisodesTotal)
				}
			}
		})
	}
}
