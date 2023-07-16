package health

import (
	"os"
	"testing"

	"github.com/lainio/err2/assert"
	"github.com/lainio/err2/try"
	"github.com/shynome/pocketbase-go-sdk/internal/pocketbase"
	"github.com/shynome/pocketbase-go-sdk/services/base"
)

var testBS *base.Service

func TestMain(m *testing.M) {
	cmd, addr := pocketbase.Start()
	defer cmd.Process.Signal(os.Interrupt)

	testBS = base.New("http://" + addr)

	m.Run()
}

func TestCheck(t *testing.T) {
	s := New(testBS)
	resp := try.To1(s.Check())
	assert.Equal(resp.Code, 200)
	t.Log(resp)
}
