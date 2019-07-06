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

type EntryService interface {
	Entry(id string) (*Entry, error)
	Entries() ([]Entry, error)
	CreateEntry(v Entry) (*Entry, error)
	UpdateEntry(v Entry) (*Entry, error)
	DeleteEntry(id string) error
}
