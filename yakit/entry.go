package yakit

const (
	INTOWN    = 0
	OUTOFTOWN = 1
)

type Entry struct {
	ID          int
	Consumption float64
	Message     string
	UsageType   int
	Vehicle     Vehicle
}
