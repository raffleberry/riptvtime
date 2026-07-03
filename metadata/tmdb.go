package metadata

import (
	"fmt"
	"os"
	"strconv"

	tmdb "github.com/cyruzin/golang-tmdb"
)

type tmdbSource struct {
	Name   string
	Url    string
	client *tmdb.Client
}

func NewTmdbSource() Source {
	client, err := tmdb.Init(os.Getenv("TMDB_API_KEY"))
	if err != nil {
		panic(err)
	}
	return &tmdbSource{
		Name:   "TMDB",
		Url:    "https://www.themoviedb.org/",
		client: client,
	}
}

func (t *tmdbSource) GetName() string {
	return t.Name
}

func (t *tmdbSource) SearchShows(query string, page int) (*SearchResultShows, error) {
	res, err := t.client.GetSearchTVShow(query, map[string]string{"page": strconv.Itoa(page)})
	if err != nil {
		return nil, err
	}
	shows := []*Show{}
	for _, s := range res.Results {
		year := int64(0)
		if len(s.FirstAirDate) >= 4 {
			year, err = strconv.ParseInt(s.FirstAirDate[0:4], 10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing year for show %s: %v, field value: %s\n", s.Name, err, s.FirstAirDate)
			}
		}
		shows = append(shows, &Show{
			ID:       s.ID,
			Title:    s.Name,
			Overview: s.Overview,
			Year:     year,
			Seasons:  0,
			Episodes: 0,
		})
	}
	return &SearchResultShows{
		Page:         res.Page,
		TotalPages:   res.TotalPages,
		TotalResults: res.TotalResults,
		Results:      shows,
	}, nil
}

func (t *tmdbSource) GetShow(id string) *Show {
	return &Show{}
}

func (t *tmdbSource) GetEpisode(id string) *Episode {
	return &Episode{}
}
