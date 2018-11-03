package yakit

type Model struct {
	ID    int
	Name  string
	Brand Brand
}

type ModelService interface {
	Model(id string) (*Model, error)
	Models(brandID string) ([]Model, error)
	CreateModel(m Model) (*Model, error)
	UpdateModel(m Model) (*Model, error)
	DeleteModel(id string) error
}
