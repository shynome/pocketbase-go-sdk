package base

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type Service struct {
	Client *resty.Client
}

func New(endpoint string) *Service {
	endpoint = strings.TrimSuffix(endpoint, "/")
	endpoint = endpoint + "/api"
	client := resty.New().
		SetBaseURL(endpoint).
		SetTimeout(10 * time.Second)
	client.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
		if !strings.HasPrefix(r.Status(), "4") {
			return nil
		}
		var msg Message
		if err := json.Unmarshal(r.Body(), &msg); err != nil {
			return err
		}
		msg.resp = r
		return &msg
	})
	return &Service{
		Client: client,
	}
}

type Message struct {
	resp    *resty.Response
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

var _ error = (*Message)(nil)

func (msg *Message) Error() string {
	return fmt.Sprintf("code: %d, message: %s", msg.Code, msg.Message)
}

func (msg *Message) Response() *resty.Response {
	return msg.resp
}
