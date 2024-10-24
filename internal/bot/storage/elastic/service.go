package elastic

import (
	"github.com/elastic/go-elasticsearch/v8"
)

type client struct {
	addresses []string
	cl        *elasticsearch.TypedClient
}

func New(addresses []string) (*client, error) {
	cl, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: addresses,
	})
	if err != nil {
		return nil, err
	}

	c := &client{cl: cl, addresses: addresses}

	return c, nil
}
