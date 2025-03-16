package record

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/lainio/err2/assert"
	"github.com/lainio/err2/try"
	"github.com/shynome/pocketbase-go-sdk/internal/pocketbase"
	"github.com/shynome/pocketbase-go-sdk/services/base"
)

var testBS *base.Service
var testUser *Service[TestUser]

func TestMain(m *testing.M) {
	cmd, addr := pocketbase.Start()
	defer cmd.Process.Signal(os.Interrupt)

	testBS = base.New("http://" + addr)
	testUser = New[TestUser](testBS, "users")

	m.Run()
}

type TestUser struct {
	UserBase
	Name string `json:"name"`
}

func TestAuth(t *testing.T) {
	resp := try.To1(testUser.AuthWithPassword("test@test.invaild", "testtest", nil))
	assert.Equal(resp.Record.Email, "test@test.invaild")
	u := try.To1(testUser.Update(resp.Record.ID, func(req *resty.Request) {
		req.SetBody(map[string]string{
			"name": "test2",
		})
	}))
	assert.Equal(u.Name, "test2")
	resp2 := try.To1(testUser.AuthRefresh(nil))
	t.Log(resp2)
}

func TestAuthWithExpired(t *testing.T) {
	resp := try.To1(testUser.AuthWithPassword("test@test.invaild", "testtest", nil))
	assert.Equal(resp.Record.Email, "test@test.invaild")

	time.Sleep(6 * time.Second)
	u := try.To1(testUser.Update(resp.Record.ID, func(req *resty.Request) {
		req.SetBody(map[string]string{
			"name": "test2",
		})
	}))
	assert.Equal(u.Name, "test2")

	collecter := NewCollectLogger()
	testUser.Client.SetLogger(collecter)

	time.Sleep(6 * time.Second)
	u = try.To1(testUser.Update(resp.Record.ID, func(req *resty.Request) {
		req.SetBody(map[string]string{
			"name": "test2",
		})
	}))
	assert.Equal(u.Name, "test2")

	assert.Equal(collecter.warn.String(), "status: 404, message: The requested resource wasn't found., Attempt 1")

	resp2 := try.To1(testUser.AuthRefresh(nil))
	t.Log(resp2)
}

type collectLogger struct {
	debug *bytes.Buffer
	err   *bytes.Buffer
	warn  *bytes.Buffer
}

var _ resty.Logger = (*collectLogger)(nil)

func NewCollectLogger() *collectLogger {
	return &collectLogger{
		debug: &bytes.Buffer{},
		err:   &bytes.Buffer{},
		warn:  &bytes.Buffer{},
	}
}

func (logger *collectLogger) Errorf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	logger.err.WriteString(s)
}
func (logger *collectLogger) Warnf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	logger.warn.WriteString(s)
}
func (logger *collectLogger) Debugf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	logger.debug.WriteString(s)
}

func TestDecodeToken(t *testing.T) {
	var s = `{"collectionId":"_pb_users_auth_","exp":1684589037,"id":"test","type":"authRecord"}`
	var jwtToken JWTToken
	try.To(json.Unmarshal([]byte(s), &jwtToken))
	t.Log(jwtToken)
}
