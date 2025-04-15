package elastic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/PAY-HERO-CONSULTING/gh-tools/apperrs"
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

type elasticConfig struct {
	client  *elasticsearch.Client
	index   string
	timeout time.Duration
}

func newNewsClient(
	ctx context.Context,
	URL string,
) (*elasticsearch.Client, error) {
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil, apperrs.NewError(err)
	}

	err = ping(context.Background(), client, URL)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func ping(ctx context.Context, client *elasticsearch.Client, url string) error {
	// Ping the Elasticsearch server to get HttpStatus, version number
	if client != nil {
		resp, err := client.Ping()
		if err != nil {
			return err
		}

		fmt.Printf("Elasticsearch returned with code %d \n", resp.StatusCode)
		return nil
	}

	return errors.New("elastic client is nil")
}
