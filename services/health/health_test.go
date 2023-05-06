package health

import (
	"testing"

	"github.com/lainio/err2/assert"
	"github.com/lainio/err2/try"
	"github.com/shynome/pocketbase-go-sdk/services/base"
)

var testBS = base.New("http://127.0.0.1:8090")

func TestCheck(t *testing.T) {
	s := New(testBS)
	resp := try.To1(s.Check())
	assert.Equal(resp.Code, 200)
	t.Log(resp)
}
