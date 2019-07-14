package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"

	"github.com/egegunes/yakit/yakit"
)

// Struct that satisfies ModelService
type ModelStore struct {
	db *sql.DB
}

func NewModelStore(db *sql.DB) *ModelStore {
	return &ModelStore{db: db}
}

// Get Model with ID
func (s ModelStore) Model(id string) (*yakit.Model, error) {
	var m yakit.Model

	stmt := `SELECT
		    models.id,
		    models.name,
		    models.type,
		    models.engine_cc,
		    models.engine_hp,
		    brands.id,
		    brands.name
		FROM models
		JOIN brands ON models.brand_id = brands.id
		WHERE models.id = $1;`

	err := s.db.QueryRow(stmt, id).Scan(
		&m.ID,
		&m.Name,
		&m.Type,
		&m.EngineCC,
		&m.EngineHP,
		&m.Brand.ID,
		&m.Brand.Name,
	)

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
				models.type,
				models.engine_cc,
				models.engine_hp,
				brands.id as brand_id,
				brands.name as brand_name
			    FROM models
			    JOIN brands
			    ON models.brand_id = brands.id`)

	if brandID != "" {
		stmt.WriteString(fmt.Sprintf(" WHERE models.brand_id = %s", brandID))
	}

	rows, err := s.db.Query(stmt.String())
	if err != nil {
		return nil, fmt.Errorf("Can't query models: %v", err)
	}

	defer rows.Close()

	var models []yakit.Model

	for rows.Next() {
		var m yakit.Model

		err = rows.Scan(
			&m.ID,
			&m.Name,
			&m.Type,
			&m.EngineCC,
			&m.EngineHP,
			&m.Brand.ID,
			&m.Brand.Name,
		)

		if err != nil {
			return nil, fmt.Errorf("Can't query models: %v", err)
		}

		models = append(models, m)
	}

	return models, nil
}

// Create a new model
func (s ModelStore) CreateModel(m yakit.Model) (*yakit.Model, error) {
	err := s.db.QueryRow(`INSERT INTO models (
		brand_id,
		name,
		type,
		engine_cc,
		engine_hp
	    ) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5
	    ) RETURNING id`,
		m.Brand.ID,
		m.Name,
		m.Type,
		m.EngineCC,
		m.EngineHP,
	).Scan(&m.ID)

	if err != nil {
		return nil, fmt.Errorf("Can't create model %d: %v", m.ID, err)
	}

	return &m, nil
}

// Update an existing model
func (s ModelStore) UpdateModel(m yakit.Model) (*yakit.Model, error) {
	_, err := s.db.Exec(`UPDATE models SET
				name = $1,
				brand_id = $2,
				type = $3,
				engine_cc = $4,
				engine_hp = $5
			    WHERE id = $6`, m.Name, m.Brand.ID, m.Type, m.EngineCC, m.EngineHP, m.ID)

	if err != nil {
		return nil, fmt.Errorf("Can't update model %d: %v", m.ID, err)
	}

	return &m, nil
}

// Delete a model
func (s ModelStore) DeleteModel(id string) error {
	_, err := s.db.Exec("DELETE FROM models WHERE id = $1", id)

	if err != nil {
		return fmt.Errorf("Can't delete model %s: %v", id, err)
	}

	return nil
}
