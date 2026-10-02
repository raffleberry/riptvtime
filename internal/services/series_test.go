package services_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/raffleberry/riptvtime/internal/config"
	"github.com/raffleberry/riptvtime/internal/db"
	"github.com/raffleberry/riptvtime/internal/meta"
	"github.com/raffleberry/riptvtime/internal/services"
	"github.com/raffleberry/riptvtime/internal/utils"
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
			nil,
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
			nil,
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
			&meta.ImdbRating{Id: "id1", Rating: 9.8, Votes: 1},
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
			nil,
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
			srv := services.NewTvService(nil, tt.db, tt.meta, tt.ipt)
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

func TestSeriesService_GetTvCacheExpireTime(t *testing.T) {
	inProdSrs := &meta.TvDetails{
		InProduction: true,
	}

	maxt := time.Now().Add(time.Hour * 61)
	mint := time.Now().Add(time.Hour * 35)
	srv := services.NewTvService(nil, nil, nil, nil)
	got := srv.GetTvCacheExpireTime(inProdSrs)
	if !(got.After(mint) && got.Before(maxt)) {
		t.Errorf("GetTvCacheExpireTime() - want [Random Time Between 36 to 60 hours] got:[%v]", got)
	}

	// got < airdate
	airDate := time.Now().Add(time.Hour * 24 * 7)
	inProdSrs.NextEpisodeToAir.AirDate = airDate
	got = srv.GetTvCacheExpireTime(inProdSrs)
	if !got.Before(airDate) {
		t.Errorf("GetTvCacheExpireTime() - want [%v] got:[%v]", airDate, got)
	}

	// < 24 hours before release
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	nextDay := today.Add(24 * time.Hour)
	inProdSrs.NextEpisodeToAir.AirDate = nextDay
	got = srv.GetTvCacheExpireTime(inProdSrs)
	if !services.IsSameDate(got, nextDay) {
		t.Errorf("GetTvCacheExpireTime() - want [%v] got:[%v]", airDate, got)
	}

	// airdate == now
	// assuming today is airdate & Next episode to air is not updated in meta service
	airDate = time.Now().UTC()
	inProdSrs.NextEpisodeToAir.AirDate = airDate
	got = srv.GetTvCacheExpireTime(inProdSrs)
	if !got.After(airDate) || got.Sub(airDate) >= time.Hour*3+time.Minute {
		t.Errorf("GetTvCacheExpireTime() - want [%v] got:[%v]", airDate, got)
	}
}

func TestSeriesService_MakeUpNext(t *testing.T) {

	fd1 := services.SeriesFullItem{
		TvDetails: &meta.TvDetails{
			NumberOfSeasons:  1,
			NumberOfEpisodes: 13,
			Seasons: []meta.TvSeason{
				meta.TvSeason{
					Name:         "Season 1",
					EpisodeCount: 13,
					SeasonNumber: 1,
				},
			},
			LastEpisodeToAir: meta.TvEpisode{
				SeasonNumber:  1,
				EpisodeNumber: 6,
			},
		},
		EpisodesAired: 6,
		EpsWatched: []services.SeriesEpisode{
			services.SeriesEpisode{S: 1, E: 1},
			services.SeriesEpisode{S: 1, E: 2},
			services.SeriesEpisode{S: 1, E: 3},
			services.SeriesEpisode{S: 1, E: 4},
			services.SeriesEpisode{S: 1, E: 5},
			services.SeriesEpisode{S: 0, E: 5},
		},
	}

	fd2, err := utils.DeepCopy[services.SeriesFullItem](fd1)
	if err != nil {
		t.Fatalf("failed to deepcopy")
	}
	// add a special episode so that airedCount = watchedCount = 6
	fd2.EpsWatched = append(fd2.EpsWatched, services.SeriesEpisode{S: 0, E: 1})

	want1 := services.SeriesEpisode{
		S: 1, E: 6,
	}

	want2 := want1

	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		cfg  *config.Config
		db   db.Db
		meta meta.Meta
		ipt  *services.ImportSvc
		// Named input parameters for target function.
		mId  int
		fd   *services.SeriesFullItem
		want *services.SeriesEpisode
	}{
		{
			"Test UpNext",
			nil,
			nil,
			nil,
			nil,
			1,
			&fd1,
			&want1,
		},
		{
			"Test UpNext with special ep watched",
			nil,
			nil,
			nil,
			nil,
			1,
			&fd2,
			&want2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := services.NewTvService(tt.cfg, tt.db, tt.meta, tt.ipt)
			got := srv.MakeUpNext(tt.mId, tt.fd)
			// TODO: update the condition below to compare got with tt.want.
			if tt.want.S != got.S || tt.want.E != got.E {
				t.Errorf("MakeUpNext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeriesService_LegitEpsWatchedCnt(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		cfg  *config.Config
		db   db.Db
		meta meta.Meta
		ipt  *services.ImportSvc
		// Named input parameters for target function.
		fd   *services.SeriesFullItem
		want int
	}{
		{
			name: "EPs watch count 2|3",
			cfg:  nil,
			db:   nil,
			meta: nil,
			ipt:  nil,
			fd: &services.SeriesFullItem{
				TvDetails: &meta.TvDetails{
					NumberOfSeasons: 1,
				},
				EpsWatched: []services.SeriesEpisode{
					{S: 1, E: 3},
					{S: 1, E: 2},
					{S: 0, E: 1},
				},
				EpisodesAired: 3,
			},
			want: 2,
		},
		{
			name: "EPs watch count 3|3",
			cfg:  nil,
			db:   nil,
			meta: nil,
			ipt:  nil,
			fd: &services.SeriesFullItem{
				TvDetails: &meta.TvDetails{
					NumberOfSeasons: 1,
				},
				EpsWatched: []services.SeriesEpisode{
					{S: 1, E: 3},
					{S: 1, E: 2},
					{S: 1, E: 1},
				},
				EpisodesAired: 3,
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := services.NewTvService(tt.cfg, tt.db, tt.meta, tt.ipt)
			got := srv.LegitEpsWatchedCnt(tt.fd)
			if got != tt.want {
				t.Errorf("LegitEpsWatchedCnt() = %v, want %v", got, tt.want)
			}
		})
	}
}
