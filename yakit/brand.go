package yakit

type Brand struct {
	ID   int
	Name string
}

type BrandService interface {
	Brand(id string) (*Brand, error)
	Brands() ([]Brand, error)
	CreateBrand(b Brand) (*Brand, error)
	UpdateBrand(b Brand) (*Brand, error)
	DeleteBrand(id string) error
}
