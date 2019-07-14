package yakit

type ModelType int

type BikeTypes struct {
	SUPERSPORT ModelType
	ENDURO     ModelType
	CRUISER    ModelType
	COMMUTER   ModelType
	SCOOTER    ModelType
	NAKED      ModelType
	CROSS      ModelType
}

var BikeType = &BikeTypes{
	SUPERSPORT: 0,
	ENDURO:     1,
	CRUISER:    2,
	COMMUTER:   3,
	SCOOTER:    4,
	NAKED:      5,
	CROSS:      6,
}

type Model struct {
	ID       int
	Name     string
	Brand    Brand
	Type     ModelType
	EngineCC int
	EngineHP int
}

type ModelService interface {
	Model(id string) (*Model, error)
	Models(brandID string) ([]Model, error)
	CreateModel(m Model) (*Model, error)
	UpdateModel(m Model) (*Model, error)
	DeleteModel(id string) error
}
