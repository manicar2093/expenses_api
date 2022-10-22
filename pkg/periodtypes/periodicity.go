package periodtypes

//go:generate stringer -type=Periodicity
type Periodicity uint64

type PeriodicityServiceImpl struct {
}

const (
	Daily Periodicity = iota + 1
	Weekly
	FourteenDays
	Paydaily
	Monthly
	BiMonthly
	FourMonthly
	SixMonthly
	Yearly
)
