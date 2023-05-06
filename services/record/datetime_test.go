package record

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/lainio/err2/assert"
	"github.com/lainio/err2/try"
)

func TestDateTime(t *testing.T) {
	now := time.Now().UTC()

	type w struct{ T DateTime }
	b := try.To1(json.Marshal(w{T: DateTime(now)}))
	var a w
	try.To(json.Unmarshal(b, &a))
	var tt = time.Time(a.T).UTC()

	a1 := now.Format(DefaultDateLayout)
	a2 := tt.Format(DefaultDateLayout)
	assert.Equal(a2, a1)

	b1 := now.Unix()
	b2 := tt.Unix()
	assert.Equal(b2, b1)
}
