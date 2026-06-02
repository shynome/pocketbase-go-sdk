package crud

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/avast/retry-go/v5"
	"github.com/tmaxmax/go-sse"
	"resty.dev/v3"
)

type SubscribeParams struct {
	url.Values
}

func (params *SubscribeParams) ToValues() url.Values {
	if params.Values == nil {
		params.Values = url.Values{}
	}
	q := params.Values
	return q
}

type Subscription[T any] struct {
	Action string `json:"action"`
	Record T      `json:"data"`
}

func (s *Service[T]) Subscribe(topic string, params *SubscribeParams, callback func(data *Subscription[T])) (_ func(), err error) {
	if params == nil {
		params = &SubscribeParams{}
	}
	q := params.ToValues()
	ctx := context.Background()
	ctx, cause := context.WithCancelCause(ctx)
	defer func() {
		if err == nil {
			return
		}
		cause(nil)
	}()
	api := s.Client.BaseURL() + "/realtime"
	wctx, connected := context.WithCancelCause(ctx)
	defer connected(nil)
	connect := func(d []byte) (err error) {
		defer func() {
			connected(err)
		}()
		var cinfo PBConnect
		if err := json.Unmarshal(d, &cinfo); err != nil {
			return err
		}
		topic := fmt.Sprintf("%s/%s", s.Collection, topic)
		body := map[string]any{
			"clientId":      cinfo.ClientId,
			"subscriptions": []string{topic},
		}
		_, err = s.Send(api, func(req *resty.Request) {
			req.SetDebug(true).
				SetMethod(http.MethodPost).
				SetQueryParamsFromValues(q).
				SetBody(body)
		})
		return err
	}
	ssec := &sse.Client{
		HTTPClient: s.Client.Client(),
		Backoff: sse.Backoff{
			MaxInterval: 5 * time.Second,
		},
	}
	retryer := retry.New(
		retry.Context(ctx),
		retry.Attempts(0),
		retry.MaxDelay(time.Second),
	)
	go retryer.Do(func() error {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, api, http.NoBody)
		if err != nil {
			return err
		}
		c := ssec.NewConnection(req)
		{
			unsub := c.SubscribeEvent("PB_CONNECT", func(e sse.Event) {
				connect([]byte(e.Data))
			})
			defer unsub()
		}
		unsub := c.SubscribeToAll(func(e sse.Event) {
			var d Subscription[T]
			err := json.Unmarshal([]byte(e.Data), &d)
			if err != nil {
				return
			}
			callback(&d)
		})
		defer unsub()
		err = c.Connect()
		return err
	})
	<-wctx.Done()
	err = context.Cause(wctx)
	if errors.Is(err, context.Canceled) {
		err = nil
	}
	if err != nil {
		return nil, err
	}
	return func() {
		cause(nil)
	}, nil
}

type PBConnect struct {
	ClientId string `json:"clientId"`
}
