package dates

import (
	"sync"
	"time"
)

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
	var (
		res  time.Time
		once = sync.Once{}
	)
	once.Do(func() {
		res = GetNormalizedDate()
	})

	return res
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
	layout := "2006-02-01T15:04:05Z"
	formated := date.Format(layout)
	t, err := time.Parse(layout, formated)
	if err != nil {
		panic(err)
	}
	return t
}

func IsMonthsAway(initTime, endTime *time.Time, months uint) bool {
	_, diffMonths, _ := diff(*initTime, *endTime)
	return uint(diffMonths) >= months
}

func IsYearsAway(initTime, endTime *time.Time, years uint) bool {
	diffYears, _, _ := diff(*initTime, *endTime)
	return uint(diffYears) >= years
}

func GetLastMonthDay(date time.Time) uint {
	firstday := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.Local)
	return uint(firstday.AddDate(0, 1, 0).Add(time.Nanosecond * -1).Day())
}

func diff(initTime, endTime time.Time) (year, month, day int) {
	if initTime.Location() != endTime.Location() {
		endTime = endTime.In(initTime.Location())
	}
	if initTime.After(endTime) {
		initTime, endTime = endTime, initTime
	}
	y1, M1, d1 := initTime.Date() //nolint:varnamelen
	y2, M2, d2 := endTime.Date()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)

	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}
