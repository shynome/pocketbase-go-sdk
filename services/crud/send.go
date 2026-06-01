package crud

import (
	"github.com/shynome/pocketbase-go-sdk/services/base"
	"resty.dev/v3"
)

func (s *Service[T]) Send(api string, initReq func(req *resty.Request)) (result base.Message, err error) {
	req := s.getReq()
	initReq(req)
	_, err = req.
		SetResult(&result).
		Execute(req.Method, api)
	return
}
