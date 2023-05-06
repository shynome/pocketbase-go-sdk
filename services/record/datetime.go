package record

// copy from https://pkg.go.dev/github.com/pocketbase/pocketbase/tools/types#DateTime

import (
	"encoding/json"
	"time"
)

const DefaultDateLayout = "2006-01-02 15:04:05.000Z"

type DateTime time.Time

func (d DateTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

func (d DateTime) String() string {
	t := time.Time(d)
	s := t.UTC().Format(DefaultDateLayout)
	return s
}

func (d *DateTime) UnmarshalJSON(b []byte) error {
	var raw string
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	t, err := time.Parse(DefaultDateLayout, raw)
	if err != nil {
		return err
	}
	*d = DateTime(t)
	return nil
}
