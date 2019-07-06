package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"

	"github.com/egegunes/yakit/yakit"
)

// Struct that satisfies EntryService
type EntryStore struct {
	db *sql.DB
}

// Create new EntryStore
func NewEntryStore(db *sql.DB) *EntryStore {
	return &EntryStore{db: db}
}

// Get Entry with ID
func (s EntryStore) Entry(id string) (*yakit.Entry, error) {
	var e yakit.Entry

	stmt := `SELECT
		entries.id,
		entries.consumption,
		entries.message,
		entries.usage_type,
		vehicles.id as vehicle_id,
		vehicles.year as vehicle_year,
		models.id as model_id,
		models.name as model_name,
		brands.id as brand_id,
		brands.name as brand_name
	    FROM entries
	    JOIN vehicles ON entries.vehicle_id = vehicles.id
	    JOIN models ON vehicles.model_id = models.id
	    JOIN brands ON models.brand_id = brands.id
	    WHERE entries.id = $1;`

	err := s.db.QueryRow(stmt, id).Scan(
		&e.ID,
		&e.Consumption,
		&e.Message,
		&e.UsageType,
		&e.Vehicle.ID,
		&e.Vehicle.Year,
		&e.Vehicle.Model.ID,
		&e.Vehicle.Model.Name,
		&e.Vehicle.Model.Brand.ID,
		&e.Vehicle.Model.Brand.Name,
	)

	if err != nil {
		return nil, fmt.Errorf("Can't query entry %s: %v", id, err)
	}

	return &e, nil
}

// List entries
func (s EntryStore) Entries() ([]yakit.Entry, error) {
	var stmt strings.Builder

	stmt.WriteString(`SELECT
			    entries.id,
			    entries.consumption,
			    entries.message,
			    entries.usage_type,
			    vehicles.id as vehicle_id,
			    vehicles.year as vehicle_year,
			    models.id as model_id,
			    models.name as model_name,
			    brands.id as brand_id,
			    brands.name as brand_name
			FROM entries
			JOIN vehicles ON entries.vehicle_id = vehicles.id
			JOIN models ON vehicles.model_id = models.id
			JOIN brands ON models.brand_id = brands.id`)

	rows, err := s.db.Query(stmt.String())

	if err != nil {
		return nil, fmt.Errorf("Can't query entries: %v", err)
	}

	defer rows.Close()

	var entries []yakit.Entry

	for rows.Next() {
		var e yakit.Entry

		err = rows.Scan(
			&e.ID,
			&e.Consumption,
			&e.Message,
			&e.UsageType,
			&e.Vehicle.ID,
			&e.Vehicle.Year,
			&e.Vehicle.Model.ID,
			&e.Vehicle.Model.Name,
			&e.Vehicle.Model.Brand.ID,
			&e.Vehicle.Model.Brand.Name,
		)

		if err != nil {
			return nil, fmt.Errorf("Can't query entries: %v", err)
		}

		entries = append(entries, e)
	}

	return entries, nil
}

// Create a new Entry
func (s EntryStore) CreateEntry(e yakit.Entry) (*yakit.Entry, error) {
	err := s.db.QueryRow(`INSERT INTO entries
				(vehicle_id,
				 consumption,
				 message,
				 usage_type)
			      VALUES ($1, $2, $3, $4) RETURNING id`,
		e.Vehicle.ID,
		e.Consumption,
		e.Message,
		e.UsageType).Scan(&e.ID)

	if err != nil {
		return nil, fmt.Errorf("Can't create entry: %v", err)
	}

	return &e, nil
}

// Update an Entry
func (s EntryStore) UpdateEntry(e yakit.Entry) (*yakit.Entry, error) {
	_, err := s.db.Exec(`UPDATE entries SET
				vehicle_id = $1,
				consumption = $2,
				message = $3,
				usage_type = $4
			    WHERE id = $5`, e.Vehicle.ID, e.Consumption, e.Message, e.UsageType, e.ID)

	if err != nil {
		return nil, fmt.Errorf("Can't update entry %d: %v", e.ID, err)
	}

	return &e, nil
}

// Delete an Entry
func (s EntryStore) DeleteEntry(id string) error {
	_, err := s.db.Exec("DELETE FROM entries WHERE id = $1", id)

	if err != nil {
		return fmt.Errorf("Can't delete entry %s: %v", id, err)
	}

	return nil
}
