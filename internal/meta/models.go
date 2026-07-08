package meta

type TvSearchResult struct {
	Id           int64
	Name         string
	Overview     string
	FirstAirDate string
	Year         int
}

type TvSearchResults struct {
	Page         int
	Results      []TvSearchResult
	TotalPages   int
	TotalResults int
}
