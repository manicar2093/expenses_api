package periodicity

//go:generate stringer -type=Periodicity
type Periodicity uint64

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
