package pocketbase_test

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/lainio/err2/assert"
	"github.com/lainio/err2/try"
	"github.com/shynome/pocketbase-go-sdk"
)

var pb = pocketbase.New("http://127.0.0.1:8090")

func TestHealth(t *testing.T) {
	resp := try.To1(pb.Health().Check())
	assert.Equal(resp.Code, 200)
}

func TestCollection(t *testing.T) {
	collection := pocketbase.NewCollection[map[string]string](pb, "public")
	r := try.To1(collection.Create(func(req *resty.Request) {
		req.SetBody(map[string]string{
			"name": "test",
		})
	}))
	defer collection.Delete(r["id"], nil)
	result := try.To1(collection.List(nil))
	assert.NotEqual(len(result.Items), 0)
}

func TestAuth(t *testing.T) {
	collection := pocketbase.NewCollection[map[string]any](pb, "users")
	try.To1(collection.AuthWithPassword("test", "testtest", nil))
	result := try.To1(collection.List(nil))
	assert.Equal(len(result.Items), 1)
}
