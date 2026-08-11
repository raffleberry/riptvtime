package services_test

import (
	"testing"
	"time"

	"github.com/raffleberry/riptvtime/internal/db"
	"github.com/raffleberry/riptvtime/internal/meta"
	"github.com/raffleberry/riptvtime/internal/services"
)

func mockDb() db.Db {
	type MockDb struct {
		db.Db
	}
	return MockDb{}
}

func TestSeriesService_MakeFeedList(t *testing.T) {

	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		c    db.Cache
		db   db.Db
		meta meta.Meta
		ipt  *services.ImportSvc
		// Named input parameters for target function.
		series          []db.TvSeries
		freshSeriesData []*services.SeriesFullItem
		want            []*services.SeriesFeedItem
	}{
		{
			"Validations - get only watching",
			nil, nil, nil, nil, []db.TvSeries{
				db.TvSeries{
					MId:            1,
					MName:          "TvMetaService",
					Name:           "Test1",
					Overview:       "Overview1",
					TrackingStatus: db.TvStatusWatching,
					Year:           2001,
				},
				db.TvSeries{
					MId:            2,
					MName:          "TvMetaService",
					Name:           "Test2",
					Overview:       "Overview2",
					TrackingStatus: db.TvStatusStopped,
					Year:           2002,
				},
			},
			[]*services.SeriesFullItem{
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
					},
					20,
					[]services.SeriesEpisode{},
				},
				&services.SeriesFullItem{
					&meta.TvDetails{
						Id:       2,
						Name:     "UpdatedTest2",
						Overview: "UpdatedOverview2",
						Year:     2012,
						LastEpisodeToAir: meta.TvEpisode{
							SeasonNumber:  1,
							EpisodeNumber: 5,
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
					},
					5,
					[]services.SeriesEpisode{},
				},
			},
			[]*services.SeriesFeedItem{
				&services.SeriesFeedItem{
					db.TvSeries{
						MId:            1,
						MName:          "TvMetaService",
						Name:           "UpdatedTest1",
						Overview:       "UpdatedOverview1",
						Year:           2011,
						TrackingStatus: db.TvStatusWatching,
					},
					1,
					1,
					0,
					2,
					7,
					false,
					"",
					time.Time{},
					time.Time{},
					time.Time{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := services.NewTvService(tt.c, tt.db, tt.meta, tt.ipt)
			got := srv.MakeFeedList(tt.series, tt.freshSeriesData)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("MakeFeedList() = %v, want %v", got, tt.want)
			}
		})
	}
}
