package crud

import (
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/shynome/pocketbase-go-sdk/services/base"
)

type Service[T any] struct {
	*base.Service
	collection string
}

func New[T any](bs *base.Service, collection string) *Service[T] {
	return &Service[T]{
		bs,
		collection,
	}
}

func (s *Service[T]) getReq() *resty.Request {
	return s.Client.R().SetPathParam("collection", s.collection)
}

type ListParams struct {
	url.Values
	Page    int
	PerPage int
	Sort    string
	Filter  string
}

func (params *ListParams) ToValues() url.Values {
	if params.Values == nil {
		params.Values = url.Values{}
	}
	q := params.Values
	if params.Page != 0 {
		q.Set("page", fmt.Sprint(params.Page))
	}
	if params.PerPage != 0 {
		q.Set("perPage", fmt.Sprint(params.PerPage))
	}
	if params.Sort != "" {
		q.Set("sort", params.Sort)
	}
	if params.Filter != "" {
		q.Set("filter", params.Filter)
	}
	return q
}

type ListResult[T any] struct {
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
	Items      []T `json:"items"`
}

func (s *Service[T]) List(params *ListParams) (result ListResult[T], err error) {
	if params == nil {
		params = &ListParams{}
	}

	q := params.ToValues()

	_, err = s.getReq().
		SetQueryParamsFromValues(q).
		SetResult(&result).
		Get("/collections/{collection}/records")
	return
}

func (s *Service[T]) Create(initBody func(req *resty.Request)) (result T, err error) {
	req := s.getReq()
	initBody(req)
	_, err = req.
		SetResult(&result).
		Post("/collections/{collection}/records")
	return
}

func (s *Service[T]) Update(id string, initBody func(req *resty.Request)) (result T, err error) {
	req := s.getReq()
	initBody(req)
	_, err = req.
		SetResult(&result).
		SetPathParam("id", id).
		Patch("/collections/{collection}/records/{id}")
	return
}

func (s *Service[T]) Delete(id string, params *url.Values) (result base.Message, err error) {
	if params == nil {
		params = &url.Values{}
	}
	_, err = s.getReq().
		SetResult(&result).
		SetQueryParamsFromValues(*params).
		SetPathParam("id", id).
		Delete("/collections/{collection}/records/{id}")
	return
}

func (s *Service[T]) One(id string, params *url.Values) (result T, err error) {
	if params == nil {
		params = &url.Values{}
	}
	_, err = s.getReq().
		SetResult(&result).
		SetQueryParamsFromValues(*params).
		SetPathParam("id", id).
		Get("/collections/{collection}/records/{id}")
	return
}

func (s *Service[T]) FullList(batch int, params *ListParams) (list []T, err error) {
	if params == nil {
		params = &ListParams{}
	}
	if batch == 0 {
		batch = 200
	}
	params.PerPage = batch
	for page := 1; true; page++ {
		params.Page = page
		var result ListResult[T]
		if result, err = s.List(params); err != nil {
			return
		}
		list = append(list, result.Items...)
		if page >= result.TotalPages {
			break
		}
	}
	return
}
