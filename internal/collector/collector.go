// Package collector is the responsible to hit RSS feeds and save into database.
package collector

type Collector interface {
	Collect(url string) ([]Item, error)
}

type Item struct {
	ID          string
	Title       string
	Description string
}
