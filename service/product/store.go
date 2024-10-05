package product

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/kapeel-mopkar/ecom/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateProduct(p types.Product) error {
	_, err := s.db.Exec("INSERT INTO products (name, description, image, price, quantity) VALUES (?, ?, ?, ?, ?)", p.Name, p.Description, p.Image, p.Price, p.Quantity)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) UpdateProduct(p types.Product) error {
	_, err := s.db.Exec("UPDATE products SET name=?, description=?, image=?, price=?, quantity=? WHERE ID=?", p.Name, p.Description, p.Image, p.Price, p.Quantity, p.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetProducts() ([]types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	products := make([]types.Product, 0)
	for rows.Next() {
		p, err := scanRowIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}
	return products, nil
}

func (s *Store) GetProductsByIDs(productIds []int) ([]types.Product, error) {
	placeholders := strings.Repeat(",?", len(productIds)-1)
	query := fmt.Sprintf("SELECT * FROM products where id in (?%s)", placeholders)

	args := make([]interface{}, len(productIds))
	for i, v := range productIds {
		args[i] = v
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	products := make([]types.Product, 0)
	for rows.Next() {
		p, err := scanRowIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}
	return products, nil
}

func scanRowIntoProduct(rows *sql.Rows) (*types.Product, error) {
	product := new(types.Product)
	err := rows.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Image,
		&product.Price,
		&product.Quantity,
		&product.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return product, nil
}
