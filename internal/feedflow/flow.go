package feedflow

import (
	"context"
	"rss-reader/internal/feedstate"
)

func Collect(ctx context.Context, producer feedstate.Producer, values []string, reject int) ([]string, error) {
	items, errs := producer.Start(values, reject)
	result := []string{}
	for value := range items {
		result = append(result, value)
	}
	if err, ok := <-errs; ok {
		return result, err
	}
	return result, nil
}
