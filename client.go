package pocketbase

import (
	"github.com/shynome/pocketbase-go-sdk/services/base"
	"github.com/shynome/pocketbase-go-sdk/services/health"
	"github.com/shynome/pocketbase-go-sdk/services/record"
)

type Client struct {
	Endpoint string
}

func New(endpoint string) *Client {
	return &Client{
		Endpoint: endpoint,
	}
}

func (c *Client) Health() *health.Service {
	bs := base.New(c.Endpoint)
	return health.New(bs)
}

func NewCollection[T any](c *Client, collection string) *record.Service[T] {
	bs := base.New(c.Endpoint)
	return record.New[T](bs, collection)
}
