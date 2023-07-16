package crud

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/lainio/err2/assert"
	"github.com/lainio/err2/try"
	"github.com/shynome/pocketbase-go-sdk/internal/pocketbase"
	"github.com/shynome/pocketbase-go-sdk/services/base"
)

var testBS *base.Service
var testPublic *Service[map[string]string]

func TestMain(m *testing.M) {
	cmd, addr := pocketbase.Start()
	defer cmd.Process.Signal(os.Interrupt)

	testBS = base.New("http://" + addr)
	testPublic = New[map[string]string](testBS, "public")

	m.Run()
}

func TestPublic(t *testing.T) {
	r, err := testPublic.Create(func(req *resty.Request) {
		req.SetBody(map[string]string{
			"name": "hello",
		})
	})
	try.To(err)

	filter := fmt.Sprintf("id='%s'", r["id"])
	l := try.To1(testPublic.List(&ListParams{Filter: filter}))
	assert.Equal(len(l.Items), 1)
	assert.Equal(l.Items[0]["name"], r["name"])

	r, err = testPublic.Update(r["id"], func(req *resty.Request) {
		req.SetBody(map[string]string{
			"name": "hello2",
		})
	})
	try.To(err)
	assert.Equal(r["name"], "hello2")

	m := try.To1(testPublic.Delete(r["id"], nil))
	t.Log(m)

	_, err = testPublic.One(r["id"], nil)
	notFound := err.(*base.Message)
	assert.Equal(notFound.Code, 404)
}

func TestFullList(t *testing.T) {
	testItems := []map[string]string{
		{"name": "a"},
		{"name": "b"},
		{"name": "c"},
	}
	items := []string{}
	defer func() {
		for _, id := range items {
			try.To1(testPublic.Delete(id, nil))
		}
	}()
	for _, item := range testItems {
		r := try.To1(testPublic.Create(func(req *resty.Request) {
			req.SetBody(item)
		}))
		items = append(items, r["id"])
	}
	filter := fmt.Sprintf("id='%s' || id='%s' || id='%s'", items[0], items[1], items[2])
	ll := try.To1(testPublic.FullList(1, &ListParams{
		Filter: filter,
		Sort:   "+created",
	}))
	for i, item := range ll {
		assert.Equal(item["id"], items[i])
	}
}
