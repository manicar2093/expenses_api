package dates

import "time"

type (
	TimeGetable interface {
		GetCurrentTime() time.Time
	}
	TimeGetter struct{}
)

func (c *TimeGetter) GetCurrentTime() time.Time {
	return GetNormalizedDate()
}

func GetNormalizedDate() time.Time {
	layout := "2006-02-01T15:04:05"
	formated := time.Now().Format(layout)
	t, err := time.Parse(layout, formated)
	if err != nil {
		panic(err)
	}
	return t
}
