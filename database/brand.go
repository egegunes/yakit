package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/egegunes/yakit/yakit"
)

// Struct that satisfies BrandService
type BrandStore struct {
	DB *sql.DB
}

// Get brand with ID
func (s BrandStore) Brand(id string) (*yakit.Brand, error) {
	var b yakit.Brand

	err := s.DB.QueryRow("SELECT id, name FROM brands WHERE id = $1", id).Scan(&b.ID, &b.Name)

	if err != nil {
		return nil, fmt.Errorf("Can't query brand %s: %v", id, err)
	}

	return &b, nil
}

// Get all brands
func (s BrandStore) Brands() ([]yakit.Brand, error) {
	rows, err := s.DB.Query("SELECT id, name FROM brands")
	defer rows.Close()

	if err != nil {
		return nil, fmt.Errorf("Can't query brands: %v", err)
	}

	var brands []yakit.Brand

	for rows.Next() {
		var b yakit.Brand

		err = rows.Scan(&b.ID, &b.Name)

		if err != nil {
			return nil, fmt.Errorf("Can't query brands: %v", err)
		}

		brands = append(brands, b)
	}

	return brands, nil
}

// Create a new brand
func (s BrandStore) CreateBrand(b yakit.Brand) (*yakit.Brand, error) {
	err := s.DB.QueryRow("INSERT INTO brands (name) VALUES ($1) RETURNING id", b.Name).Scan(&b.ID)

	if err != nil {
		return nil, fmt.Errorf("Can't create brand %d: %v", b.ID, err)
	}

	return &b, nil
}

// Update an existing brand
func (s BrandStore) UpdateBrand(b yakit.Brand) (*yakit.Brand, error) {
	_, err := s.DB.Exec("UPDATE brands SET name = $1 WHERE id = $2", b.Name, b.ID)

	if err != nil {
		return nil, fmt.Errorf("Can't update brand %d: %v", b.ID, err)
	}

	return &b, nil
}

// Delete a brand
func (s BrandStore) DeleteBrand(id string) error {
	_, err := s.DB.Exec("DELETE FROM brands WHERE id=$1", id)

	if err != nil {
		return fmt.Errorf("Can't delete brand %s: %v", id, err)
	}

	return nil
}
