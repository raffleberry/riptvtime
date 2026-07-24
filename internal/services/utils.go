package services

import (
	"strings"

	"gitlab.com/raffleberry/riptvtime/internal/db"
	"gitlab.com/raffleberry/riptvtime/internal/meta"
)

func MetaToDbSeason(mId int, mSd *meta.TvSeason) *db.TvSeason {
	dbEps := []db.TvEpisode{}
	for _, e := range mSd.Episodes {
		dbEps = append(dbEps, db.TvEpisode{
			MId:       int64(e.Id),
			MName:     e.MName,
			Name:      e.Name,
			SeriesMId: int64(mId),
			Overview:  e.Overview,
			Season:    e.SeasonNumber,
			Episode:   e.EpisodeNumber,
			Runtime:   e.Runtime,
			AirDate:   e.AirDate,
		})
	}

	dbSd := db.TvSeason{
		MId:       int64(mSd.Id),
		MName:     mSd.MName,
		AirDate:   mSd.AirDate,
		SeriesMId: int64(mId),
		Season:    mSd.SeasonNumber,
		Name:      mSd.Name,
		Overview:  mSd.Overview,
		Episodes:  dbEps,
	}
	return &dbSd
}

func DbToMetaSeason(dbSd *db.TvSeason) *meta.TvSeason {
	mSd := meta.TvSeason{
		Id:           int(dbSd.MId),
		Name:         dbSd.Name,
		Overview:     dbSd.Overview,
		SeasonNumber: dbSd.Season,
		AirDate:      dbSd.AirDate,
		MName:        dbSd.MName,
	}

	for _, e := range dbSd.Episodes {
		mSd.Episodes = append(mSd.Episodes, meta.TvEpisode{
			Id:            int(e.MId),
			Name:          e.Name,
			Overview:      e.Overview,
			SeasonNumber:  e.Season,
			EpisodeNumber: e.Episode,
			Runtime:       e.Runtime,
			AirDate:       e.AirDate,
			MName:         e.MName,
			Year:          e.AirDate.Year(),
		})
	}

	mSd.EpisodeCount = len(mSd.Episodes)

	return &mSd
}

func isLegitSeason(sName string, sNo int, sCnt int) bool {
	if sNo < 1 || sNo > sCnt {
		return false
	}
	if sName == "" || !strings.Contains(strings.ToLower(sName), "special") {
		return false
	}
	return true
}
