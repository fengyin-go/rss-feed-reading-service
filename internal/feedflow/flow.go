package feedflow

import (
	"context"
	"rss-reader/internal/feedstate"
)

func Collect(ctx context.Context, producer feedstate.Producer, values []string, reject int) ([]string, error) {
	items, _ := producer.Start(values, reject)
	result := []string{}
	for value := range items {
		result = append(result, value)
	}
	return result, nil
}
