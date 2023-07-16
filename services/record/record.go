package record

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/shynome/pocketbase-go-sdk/services/base"
	"github.com/shynome/pocketbase-go-sdk/services/crud"
)

type Service[T any] struct {
	*crud.Service[T]

	locker sync.Locker
	login  func() error
}

func New[T any](bs *base.Service, collection string) (s *Service[T]) {
	s = &Service[T]{
		Service: crud.New[T](bs, collection),
		locker:  &sync.Mutex{},
	}

	s.Client.SetRetryCount(1)
	s.Client.AddRetryCondition(func(r *resty.Response, err error) bool {
		if err == nil {
			return false
		}
		if r == nil || r.Request == nil {
			return false
		}
		switch r.StatusCode() {
		case 401: // token 已过期
		case 404: // token 过期也会触发 404
		default:
			return false
		}

		s.locker.Lock()
		defer s.locker.Unlock()

		token := s.Client.Token
		if token == "" {
			return false
		}

		if s.login == nil {
			return false
		}

		tokenArr := strings.Split(token, ".")
		if len(tokenArr) < 2 {
			return false
		}
		b, err := base64.RawStdEncoding.DecodeString(tokenArr[1])
		if err != nil {
			return false
		}
		var jwtToken JWTToken
		if err := json.Unmarshal(b, &jwtToken); err != nil {
			return false
		}
		now := time.Now().Unix()
		if d := jwtToken.Exp - now; d > 0 { // token 没有过期的话不要重试
			return false
		}

		if err := s.login(); err != nil { // login 失败的话不要重试
			return false
		}
		return true
	})

	return s
}

type JWTToken struct {
	CollectionID string `json:"collectionId"`
	Exp          int64  `json:"exp"`
	ID           string `json:"id"`
	Type         string `json:"type"`
}

func (s *Service[T]) getReq() *resty.Request {
	return s.Client.R()
}

type AuthProvider struct {
	Name                string `json:"name"`
	State               string `json:"state"`
	CodeVerifier        string `json:"codeVerifier"`
	CodeChallenge       string `json:"codeChallenge"`
	CodeChallengeMethod string `json:"codeChallengeMethod"`
	AuthUrl             string `json:"authUrl"`
}

type AuthMethods struct {
	UsernamePassword bool           `json:"usernamePassword"`
	EmailPassword    bool           `json:"emailPassword"`
	AuthProviders    []AuthProvider `json:"authProviders"`
}

func (s *Service[T]) ListAuthMethods(params *url.Values) (result AuthMethods, err error) {
	if params == nil {
		params = &url.Values{}
	}
	_, err = s.getReq().
		SetQueryParamsFromValues(*params).
		SetResult(&result).
		Get(fmt.Sprintf("/collections/%s/auth-methods", s.Collection))
	return
}

type AuthResponse[T any] struct {
	Record T      `json:"record"`
	Token  string `json:"token"`
}

type RecordQueryParams struct {
	url.Values
	Expand string
	Body   map[string]any
}

func (params *RecordQueryParams) ToValues() url.Values {
	if params.Values == nil {
		params.Values = url.Values{}
	}
	q := params.Values
	q.Set("expand", params.Expand)
	return q
}

func (s *Service[T]) AuthWithPassword(identity, password string, params *RecordQueryParams) (result AuthResponse[T], err error) {
	if params == nil {
		params = &RecordQueryParams{}
	}
	if params.Body == nil {
		params.Body = make(map[string]any)
	}
	body := params.Body

	body["identity"] = identity
	body["password"] = password

	_, err = s.getReq().
		SetBody(body).
		SetQueryParamsFromValues(params.ToValues()).
		SetResult(&result).
		Post(fmt.Sprintf("/collections/%s/auth-with-password", s.Collection))
	if err != nil {
		return
	}

	s.Client.SetAuthToken(result.Token)
	s.login = func() error {
		_, err := s.AuthWithPassword(identity, password, nil)
		return err
	}

	return
}

func (s *Service[T]) AuthRefresh(params *RecordQueryParams) (result AuthResponse[T], err error) {
	if params == nil {
		params = &RecordQueryParams{}
	}

	_, err = s.getReq().
		SetQueryParamsFromValues(params.ToValues()).
		SetResult(&result).
		Post(fmt.Sprintf("/collections/%s/auth-refresh", s.Collection))
	if err != nil {
		return
	}

	s.Client.SetAuthToken(result.Token)
	return
}
