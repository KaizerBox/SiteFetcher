package bestbuy

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadAvailabilityQueryConfig(path string) (map[string]string, error) {
	if path == "" {
		path = "internal/fetcher/bestbuy/availability_query_config.json"
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read availability query config: %w", err)
	}

	queryConfig := map[string]string{}
	if err := json.Unmarshal(content, &queryConfig); err != nil {
		return nil, fmt.Errorf("decode availability query config: %w", err)
	}

	return queryConfig, nil
}
