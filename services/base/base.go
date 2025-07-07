package base

import (
	"fmt"
	"strings"
	"time"

	"resty.dev/v3"
)

type Service struct {
	Client *resty.Client
}

func New(endpoint string) *Service {
	endpoint = strings.TrimSuffix(endpoint, "/")
	endpoint = endpoint + "/api"
	client := resty.New().
		SetBaseURL(endpoint).
		SetTimeout(10 * time.Second).
		SetError(&Message{})
	client.AddResponseMiddleware(func(c *resty.Client, r *resty.Response) error {
		if !r.IsError() {
			return nil
		}
		if msg, ok := r.Error().(*Message); ok {
			msg.resp = r
			return msg
		}
		var msg = Message{
			resp:    r,
			Status:  r.StatusCode(),
			Code:    r.StatusCode(),
			Message: r.String(),
		}
		return &msg
	})
	return &Service{
		Client: client,
	}
}

type Message struct {
	resp    *resty.Response
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`

	// Deprecated: after v0.23.0 replaced with status
	Code int `json:"code"`
}

var _ error = (*Message)(nil)

func (msg *Message) Error() string {
	if msg.Status != 0 {
		return fmt.Sprintf("status: %d, message: %s", msg.Status, msg.Message)
	}
	return fmt.Sprintf("code: %d, message: %s", msg.Code, msg.Message)
}

func (msg *Message) Response() *resty.Response {
	return msg.resp
}
