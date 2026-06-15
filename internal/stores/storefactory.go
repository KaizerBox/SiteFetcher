package stores

import (
	"fmt"
	"net/http"

	"siteFetcher/internal/fetcher"
	"siteFetcher/internal/fetcher/bestbuy"
)

// NewStore creates a store implementation by kind. Currently supports "bestbuy".
func NewStore(kind string, client *http.Client) (fetcher.Store, error) {
	switch kind {
	case "bestbuy":
		return bestbuy.NewBestBuyStore(client)
	default:
		return nil, fmt.Errorf("unsupported store type: %s", kind)
	}
}
