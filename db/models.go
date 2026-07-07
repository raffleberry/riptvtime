package db

import (
	"time"

	"gorm.io/gorm"
)

type TvStatus int

const (
	TvStatusWatching TvStatus = iota
	TvStatusStopped
	TvStatusCompleted
)

func (s TvStatus) String() string {
	strings := [...]string{"not watching", "watching", "stopped", "completed"}
	if s < TvStatusWatching || int(s) >= len(strings) {
		return "unknown"
	}
	return strings[s]
}

type TvSeries struct {
	gorm.Model
	TmdbId         int64
	Name           string
	Overview       string
	Genres         string
	Year           int
	FirstAirDate   time.Time
	TrackingStatus TvStatus
	RuntimeApprox  int
}

type TvTracking struct {
	gorm.Model
	EpisodeTmdbId int64
	SeriesTmdbId  int64
	Name          string
	Overview      string
	Season        int
	Episode       int
	Runtime       int
}

type TvEpisode struct {
	gorm.Model
	TmdbID     int64
	SeriesTmdb int64
	Name       string
	Overview   string
	Season     int
	Episode    int
	Runtime    int
	AirDate    time.Time
}

type TvSeason struct {
	gorm.Model
	TmdbID     int64
	AirDate    time.Time
	SeriesTmdb int64
	Season     int
	Name       string
	Overview   string
	Episodes   []TvEpisode `gorm:"foreignkey:SeriesTmdb;references:SeriesTmdb"`
}

// func GetSeries(tmdbId int64) (*TvSeries, error) {
// 	r, err := Conn.Query(`SELECT
// 		Id, TmdbID, name,
// 		 overview, genres, year,
// 		 seasons, episodes, first_air_date,
// 		 last_air_date, in_production, tracking_status
// 		FROM tv_series
// 		WHERE tmdb_id = ?`, tmdbId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var t TvSeries
// 	for r.Next() {
// 		if err := r.Scan(&t.Id, &t.TmdbID, &t.Name, &t.Overview, &t.Genres, &t.Year, &t.Seasons, &t.Episodes, &t.FirstAirDate, &t.LastAirDate, &t.InProduction, &t.TrackingStatus); err != nil {
// 			return nil, err
// 		}
// 	}

// 	return &t, nil
// }

// func (s *TvSeries) Save() (int64, error) {
// 	series, err := GetSeries(s.TmdbID)
// 	if err != nil {
// 		return 0, err
// 	}
// 	if series != nil {
// 		slog.Debug("tv series already exists, updating it", "tmdb_id", s.TmdbID)
// 		res, err := Conn.Exec(`
// 			UPDATE tv_series
// 			SET name = ?, overview = ?, genres = ?,
// 			year = ?, seasons = ?, episodes = ?,
// 			first_air_date = ?, last_air_date = ?, in_production = ?,
// 			tracking_status = ? WHERE tmdb_id = ?
// 		`, s.Name, s.Overview, s.Genres,
// 			s.Year, s.Seasons, s.Episodes,
// 			s.FirstAirDate, s.LastAirDate, s.InProduction,
// 			s.TrackingStatus, s.TmdbID)

// 		if err != nil {
// 			return series.Id, err
// 		}

// 		rowsCnt, err := res.RowsAffected()
// 		if err != nil {
// 			return series.Id, err
// 		}
// 		slog.Debug("tvseries updated", "rows affected", rowsCnt)
// 		return series.Id, nil
// 	}

// 	if s.TrackingStatus == TvStatusNotWatching {
// 		s.TrackingStatus = TvStatusWatching
// 	}

// 	res, err := Conn.Exec(
// 		`INSERT INTO tv_series (tmdb_id, name, overview,
// 								genres, year, seasons,
// 								episodes, first_air_date, last_air_date,
// 								in_production, tracking_status
// 		) VALUES (?, ?, ?,   ?, ?, ?,   ?, ?, ?,   ?, ?)`,
// 		s.TmdbID, s.Name, s.Overview,
// 		s.Genres, s.Year, s.Seasons,
// 		s.Episodes, s.FirstAirDate, s.LastAirDate,
// 		s.InProduction, s.TrackingStatus)

// 	if err != nil {
// 		return -1, err
// 	}

// 	insertId, err := res.LastInsertId()

// 	if err != nil {
// 		return -1, err
// 	}

// 	slog.Debug("Row inserted", "id", insertId)

// 	s.Id = insertId

// 	return insertId, nil

// }

// func NewTvShowFromTmdb(t *tmdb.TVDetails) *TvSeries {

// 	var genres strings.Builder
// 	for i, g := range t.Genres {
// 		if i != 0 {
// 			genres.WriteString(",")
// 		}
// 		genres.WriteString(g.Name)
// 	}

// 	var year int
// 	if len(t.FirstAirDate) >= 4 {
// 		var err error
// 		year, err = strconv.Atoi(t.FirstAirDate[:4])
// 		if err != nil {
// 			log.Printf("Error parsing year for show %s: %v, field value: %s\n", t.Name, err, t.FirstAirDate)
// 			year = 0
// 		}
// 	}

// 	firstAirDate, err := time.Parse(time.DateOnly, t.FirstAirDate)
// 	if err != nil {
// 		log.Printf("Error parsing firstAirDate for show %s: %v, field value: %s\n", t.Name, err, t.FirstAirDate)
// 	}

// 	lastAirDate, err := time.Parse(time.DateOnly, t.LastAirDate)
// 	if err != nil {
// 		log.Printf("Error parsing firstAirDate for show %s: %v, field value: %s\n", t.Name, err, t.FirstAirDate)
// 	}

// 	series := &TvSeries{
// 		TmdbID:         t.ID,
// 		Name:           t.Name,
// 		Overview:       t.Overview,
// 		Genres:         genres.String(),
// 		Year:           year,
// 		Seasons:        t.NumberOfSeasons,
// 		Episodes:       t.NumberOfEpisodes,
// 		FirstAirDate:   firstAirDate,
// 		LastAirDate:    lastAirDate,
// 		InProduction:   t.InProduction,
// 		TrackingStatus: TvStatusNotWatching,
// 	}

// 	alreadySeries, err := GetSeries(t.ID)
// 	if err != nil {
// 		slog.Warn("Error getting series status", "tmdb_id", t.ID, "err", err)
// 	}

// 	if alreadySeries != nil {
// 		series.TrackingStatus = alreadySeries.TrackingStatus
// 	}

// 	return series
// }

// func NewTvEpisodeFromTmdb(seriesId int, ep *tmdb.TVEpisodeDetails) *TvEpisode {

// 	airDate, err := time.Parse(time.DateOnly, ep.AirDate)
// 	if err != nil {
// 		log.Printf("Error parsing AirDate for show %v, episode %s: %v, field value: %v\n", seriesId, ep.Name, err, ep.AirDate)
// 	}

// 	return &TvEpisode{
// 		Id:       -1,
// 		TmdbID:   ep.ID,
// 		Name:     ep.Name,
// 		SeriesId: seriesId,
// 		Overview: ep.Overview,
// 		AirDate:  airDate,
// 		Season:   ep.SeasonNumber,
// 		Episode:  ep.EpisodeNumber,
// 		Runtime:  ep.Runtime,
// 	}
// }
