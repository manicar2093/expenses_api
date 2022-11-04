package periodtypes

import (
	"fmt"
	"time"

	"github.com/manicar2093/expenses_api/pkg/dates"
)

//go:generate stringer -type=Periodicity
type Periodicity uint64

const (
	Empty Periodicity = iota
	Daily
	Weekly
	FourteenDays
	Paydaily
	Monthly
	BiMonthly
	FourMonthly
	SixMonthly
	Yearly
)

func (c Periodicity) GetTimeQuantity() int {
	var time int
	switch c {
	case BiMonthly:
		time = 2
	case FourMonthly:
		time = 4
	case SixMonthly:
		time = 6
	case Monthly,
		Yearly:
		time = 1
	default:
		panic(fmt.Sprintf("periodicity's time '%v' not registered", c))
	}
	return time
}

func (c Periodicity) GetExpensesQuantity(date time.Time) uint {
	var time uint
	switch c {
	case Daily:
		time = uint(dates.GetLastMonthDay(date))
	case Weekly:
		time = 4
	case FourteenDays,
		Paydaily:
		time = 2
	case Empty,
		Monthly,
		BiMonthly,
		FourMonthly,
		SixMonthly,
		Yearly:
		time = 1
	default:
		panic(fmt.Sprintf("periodicity's expenses counter '%v' not registered", c))
	}
	return time
}

func (c Periodicity) GetTimeValidator() (time func(*time.Time, *time.Time, uint) bool) {
	switch c {
	case Monthly,
		BiMonthly,
		FourMonthly,
		SixMonthly:
		time = dates.IsMonthsAway
	case Yearly:
		time = dates.IsYearsAway
	default:
		time = nil
	}
	return time
}
