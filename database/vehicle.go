package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"

	"github.com/egegunes/yakit/yakit"
)

// Struct that satisfies VehicleService
type VehicleStore struct {
	DB *sql.DB
}

// Get Vehicle with ID
func (s VehicleStore) Vehicle(id string) (*yakit.Vehicle, error) {
	var v yakit.Vehicle

	stmt := `SELECT
		    vehicles.id,
		    vehicles.year,
		    models.id as model_id,
		    models.name as model_name,
		    brands.id as brand_id,
		    brands.name as brand_name
		FROM vehicles
		JOIN models ON vehicles.model_id = models.id
		JOIN brands ON models.brand_id = brands.id
		WHERE models.id = $1;`

	err := s.DB.QueryRow(stmt, id).Scan(&v.ID, &v.Year, &v.Model.ID, &v.Model.Name, &v.Model.Brand.ID, &v.Model.Brand.Name)

	if err != nil {
		return nil, fmt.Errorf("Can't query vehicle %s: %v", id, err)
	}

	return &v, nil
}

// Get all vehicles
func (s VehicleStore) Vehicles(modelID string, brandID string) ([]yakit.Vehicle, error) {
	var stmt strings.Builder

	stmt.WriteString(`SELECT
			    vehicles.id,
			    vehicles.year,
			    models.id as model_id,
			    models.name as model_name,
			    brands.id as brand_id,
			    brands.name as brand_name
			FROM vehicles
			JOIN models ON vehicles.model_id = models.id
			JOIN brands ON models.brand_id = brands.id`)

	if modelID != "" {
		stmt.WriteString(fmt.Sprintf(" WHERE model_id = %s", modelID))
	}

	if modelID == "" && brandID != "" {
		stmt.WriteString(fmt.Sprintf(" WHERE brand_id = %s", brandID))
	}

	rows, err := s.DB.Query(stmt.String())
	defer rows.Close()

	if err != nil {
		return nil, fmt.Errorf("Can't query vehicles: %v", err)
	}

	var vehicles []yakit.Vehicle

	for rows.Next() {
		var v yakit.Vehicle

		err = rows.Scan(&v.ID, &v.Year, &v.Model.ID, &v.Model.Name, &v.Model.Brand.ID, &v.Model.Brand.Name)

		if err != nil {
			return nil, fmt.Errorf("Can't query vehicles: %v", err)
		}

		vehicles = append(vehicles, v)
	}

	return vehicles, nil
}

// Create a new Vehicle
func (s VehicleStore) CreateVehicle(v yakit.Vehicle) (*yakit.Vehicle, error) {
	err := s.DB.QueryRow("INSERT INTO vehicles (model_id, year) VALUES ($1, $2) RETURNING id", v.Model.ID, v.Year).Scan(&v.ID)

	if err != nil {
		return nil, fmt.Errorf("Can't create vehicles %d: %v", v.ID, err)
	}

	return &v, nil
}

// Update an existing Vehicle
func (s VehicleStore) UpdateVehicle(v yakit.Vehicle) (*yakit.Vehicle, error) {
	_, err := s.DB.Exec("UPDATE vehicles SET year = $1, model_id = $2 WHERE id = $3", v.Year, v.Model.ID, v.ID)

	if err != nil {
		return nil, fmt.Errorf("Can't update vehicle %d: %v", v.ID, err)
	}

	return &v, nil
}

// Delete a Vehicle
func (s VehicleStore) DeleteVehicle(id string) error {
	_, err := s.DB.Exec("DELETE FROM vehicles WHERE id = $1", id)

	if err != nil {
		return fmt.Errorf("Can't delete Vehicle %s: %v", id, err)
	}

	return nil
}
