package health

import "github.com/shynome/pocketbase-go-sdk/services/base"

type Service struct {
	*base.Service
}

func New(bs *base.Service) *Service {
	return &Service{
		bs,
	}
}

type CheckResponse = base.Message

func (s *Service) Check() (resp CheckResponse, err error) {
	_, err = s.Client.R().SetResult(&resp).Get("/health")
	return
}
