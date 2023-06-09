package repositories

import (
	"config/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type ProductsRepository struct {
	conn *pgx.Conn
}

func NewProductsRepository(conn *pgx.Conn) *ProductsRepository {
	return &ProductsRepository{conn: conn}
}

func (repo *ProductsRepository) GetOne(productCode string) (*models.Product, error) {
	sql := "select product_code, description from products where product_code = $1"
	var product models.Product
	if err := repo.conn.QueryRow(context.Background(), sql, productCode).Scan(&product.ProductCode, &product.Description); err != nil {
		return nil, err
	}
	return &product, nil
}

func (repo *ProductsRepository) GetMany() (*[]models.Product, error) {
	sql := "select product_code, description from products"
	var products []models.Product
	rows, err := repo.conn.Query(context.Background(), sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ProductCode, &product.Description); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return &products, nil
}

func (repo *ProductsRepository) Create(product *models.Product) error {
	sql := `
insert into products (product_code, description)
values ($1, $2)
	`
	tag, err := repo.conn.Exec(context.Background(), sql, product.ProductCode, product.Description)
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("%v rows affected by insert (expected 1)", tag.RowsAffected())
	}
	return err
}

func (repo *ProductsRepository) Update(product *models.Product) error {
	sql := `
update products
set description = $2
where product_code = $1
	`
	tag, err := repo.conn.Exec(context.Background(), sql, product.ProductCode, product.Description)
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("%v rows affected by update (expected 1)", tag.RowsAffected())
	}
	return err
}

func (repo *ProductsRepository) Migrate() error {
	var currentVersion int
	row := repo.conn.QueryRow(context.Background(), "select current_version from table_versions where table_name = 'products'")
	err := row.Scan(&currentVersion)
	if err != nil {
		if err == pgx.ErrNoRows {
			currentVersion = 0
		} else {
			return err
		}
	}
	if currentVersion < 1 {
		_, err = repo.conn.Exec(context.Background(), `
create table products (
	product_code text primary key,
	description text not null
)
		`)
		if err != nil {
			return err
		}
		_, err = repo.conn.Exec(context.Background(), `
insert into table_versions (table_name, current_version)
values ('products', 1)
		`)
		if err != nil {
			return err
		}
	}
	return nil
}
