package dates

import "time"

type (
	TimeGetable interface {
		GetCurrentTime() time.Time
		GetNextMonthAtFirtsDay() time.Time
	}
	TimeMonthModifier interface {
	}
	TimeGetter struct{}
)

func (c *TimeGetter) GetCurrentTime() time.Time {
	return GetNormalizedDate()
}

func (c *TimeGetter) GetNextMonthAtFirtsDay() time.Time {
	date := GetNormalizedDate()
	return NormalizeDate(
		time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, time.Local),
	)
}

func GetNormalizedDate() time.Time {
	return NormalizeDate(time.Now())
}

func NormalizeDate(date time.Time) time.Time {
	layout := "2006-02-01T15:04:05"
	formated := date.Format(layout)
	t, err := time.Parse(layout, formated)
	if err != nil {
		panic(err)
	}
	return t
}
