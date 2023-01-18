package period

const (
	Daily         Periodicity = "Daily"
	Weekly        Periodicity = "Weekly"
	FourteenDaily Periodicity = "FourteenDaily"
	Paydaily      Periodicity = "Paydaily"
	Monthly       Periodicity = "Monthly"
	BiMonthly     Periodicity = "BiMonthly"
	FourMonthly   Periodicity = "FourMonthly"
	SixMonthly    Periodicity = "SixMonthly"
	Yearly        Periodicity = "Yearly"
)

var Periods = Periodicities{Daily, Weekly, FourteenDaily, Paydaily, Monthly, BiMonthly, FourMonthly, SixMonthly, Yearly}

type (
	Periodicity   string
	Periodicities []Periodicity
)
