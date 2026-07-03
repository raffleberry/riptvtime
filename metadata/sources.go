package metadata

type Source interface {
	GetName() string
	SearchShows(query string, page int) (*SearchResultShows, error)
	GetShow(id string) *Show
	GetEpisode(id string) *Episode
}

type Show struct {
	ID       int64
	Title    string
	Overview string
	Year     int64
	Seasons  int64
	Episodes int64
}

type Episode struct {
	ID          string
	ShowID      string
	Title       string
	Description string
	Season      int64
	Episode     int64
}

type SearchResultShows struct {
	Page         int64
	TotalPages   int64
	TotalResults int64
	Results      []*Show
}
