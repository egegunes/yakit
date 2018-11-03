package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"

	"yakit/yakit"
)

// Struct that satisfies ModelService
type ModelStore struct {
	DB *sql.DB
}

// Get Model with ID
func (s ModelStore) Model(id string) (*yakit.Model, error) {
	var m yakit.Model

	stmt := `SELECT models.id, models.name, brands.id, brands.name FROM models
		JOIN brands ON models.brand_id = brands.id
		WHERE models.id = $1;`

	err := s.DB.QueryRow(stmt, id).Scan(&m.ID, &m.Name, &m.Brand.ID, &m.Brand.Name)

	if err != nil {
		return nil, fmt.Errorf("Can't query model %s: %v", id, err)
	}

	return &m, nil
}

// Get all models
func (s ModelStore) Models(brandID string) ([]yakit.Model, error) {
	var stmt strings.Builder

	stmt.WriteString(`SELECT
				models.id,
				models.name,
				brands.id as brand_id,
				brands.name as brand_name
			    FROM models
			    JOIN brands
			    ON models.brand_id = brands.id`)

	if brandID != "" {
		stmt.WriteString(fmt.Sprintf(" WHERE models.brand_id = %s", brandID))
	}

	rows, err := s.DB.Query(stmt.String())
	defer rows.Close()

	if err != nil {
		return nil, fmt.Errorf("Can't query models: %v", err)
	}

	var models []yakit.Model

	for rows.Next() {
		var m yakit.Model

		err = rows.Scan(&m.ID, &m.Name, &m.Brand.ID, &m.Brand.Name)

		if err != nil {
			return nil, fmt.Errorf("Can't query models: %v", err)
		}

		models = append(models, m)
	}

	return models, nil
}

// Create a new model
func (s ModelStore) CreateModel(m yakit.Model) (*yakit.Model, error) {
	err := s.DB.QueryRow("INSERT INTO models (brand_id, name) VALUES ($1, $2) RETURNING id", m.Brand.ID, m.Name).Scan(&m.ID)

	if err != nil {
		return nil, fmt.Errorf("Can't create model %d: %v", m.ID, err)
	}

	return &m, nil
}

// Update an existing model
func (s ModelStore) UpdateModel(m yakit.Model) (*yakit.Model, error) {
	_, err := s.DB.Exec("UPDATE models SET name = $1, brand_id = $2 WHERE id = $3", m.Name, m.Brand.ID, m.ID)

	if err != nil {
		return nil, fmt.Errorf("Can't update model %d: %v", m.ID, err)
	}

	return &m, nil
}

// Delete a model
func (s ModelStore) DeleteModel(id string) error {
	_, err := s.DB.Exec("DELETE FROM models WHERE id = $1", id)

	if err != nil {
		return fmt.Errorf("Can't delete model %s: %v", id, err)
	}

	return nil
}
