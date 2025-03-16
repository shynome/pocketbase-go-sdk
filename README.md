# Description

This is a go client sdk for pocketbase, not implement subscribe (sse)

if you want to use this, you can look the [client_test.go](./client_test.go) file

```go
package pocketbase_test

import (
	"os"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/lainio/err2/assert"
	"github.com/lainio/err2/try"
	"github.com/shynome/pocketbase-go-sdk"
	test "github.com/shynome/pocketbase-go-sdk/internal/pocketbase"
)

var pb *pocketbase.Client

func TestMain(m *testing.M) {
	cmd, addr := test.Start()
	defer cmd.Process.Signal(os.Interrupt)

	pb = pocketbase.New("http://" + addr)

	m.Run()
}

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
	defer func() {
		_, err := collection.Delete(r["id"], nil)
		if err != nil {
			t.Error(err)
		}
	}()
	result := try.To1(collection.List(nil))
	assert.NotEqual(len(result.Items), 0)
}

func TestAuth(t *testing.T) {
	collection := pocketbase.NewCollection[map[string]any](pb, "users")
	try.To1(collection.AuthWithPassword("test@test.invaild", "testtest", nil))
	result := try.To1(collection.List(nil))
	assert.Equal(len(result.Items), 1)
}

```