package crud

import (
	"testing"
	"time"
)

func TestSubscribe(t *testing.T) {
	unsub, err := testPublic.Subscribe("*", nil, func(data *Subscription[map[string]string]) {})
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(6 * time.Minute)
	unsub()
	time.Sleep(30 * time.Second)
}
