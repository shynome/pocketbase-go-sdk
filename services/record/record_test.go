package record

import (
	"encoding/json"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/lainio/err2/assert"
	"github.com/lainio/err2/try"
	"github.com/shynome/pocketbase-go-sdk/services/base"
)

var testBS = base.New("http://127.0.0.1:8090")
var testUser = New[TestUser](testBS, "users")

type TestUser struct {
	UserBase
	Name string `json:"name"`
}

func TestAuth(t *testing.T) {
	resp := try.To1(testUser.AuthWithPassword("test", "testtest", nil))
	assert.Equal(resp.Record.Username, "test")
	u := try.To1(testUser.Update(resp.Record.ID, func(req *resty.Request) {
		req.SetBody(map[string]string{
			"name": "test2",
		})
	}))
	assert.Equal(u.Name, "test2")
	resp2 := try.To1(testUser.AuthRefresh(nil))
	t.Log(resp2)
}

func TestDecodeToken(t *testing.T) {
	var s = `{"collectionId":"_pb_users_auth_","exp":1684589037,"id":"test","type":"authRecord"}`
	var jwtToken JWTToken
	try.To(json.Unmarshal([]byte(s), &jwtToken))
	t.Log(jwtToken)
}
