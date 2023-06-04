package pocketbase

import (
	"github.com/shynome/pocketbase-go-sdk/services/base"
	"github.com/shynome/pocketbase-go-sdk/services/health"
	"github.com/shynome/pocketbase-go-sdk/services/record"
)

type Client struct {
	Endpoint string

	authStore *base.Service
}

func New(endpoint string) *Client {
	return &Client{
		Endpoint: endpoint,
	}
}

func (c *Client) Health() *health.Service {
	bs := c.getBS()
	return health.New(bs)
}

func (c *Client) SetAuthStore(r *record.Service[record.UserBase]) {
	c.authStore = r.Service.Service
}

func (c *Client) getBS() *base.Service {
	if c.authStore != nil {
		return c.authStore
	}
	return base.New(c.Endpoint)
}

func NewCollection[T any](c *Client, collection string) *record.Service[T] {
	bs := c.getBS()
	return record.New[T](bs, collection)
}

func NewAdmin[T any](c *Client) *record.Service[T] {
	bs := c.getBS()
	return record.New[T](bs, "-admins")
}
