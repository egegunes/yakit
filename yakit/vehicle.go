package yakit

type Vehicle struct {
	ID    int
	Year  int
	Model Model
}

type VehicleService interface {
	Vehicle(id string) (*Vehicle, error)
	Vehicles(modelID string, brandID string) ([]Vehicle, error)
	CreateVehicle(v Vehicle) (*Vehicle, error)
	UpdateVehicle(v Vehicle) (*Vehicle, error)
	DeleteVehicle(id string) error
}
