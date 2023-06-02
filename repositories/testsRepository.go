package repositories

import (
	"config/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type TestsRepository struct {
	conn *pgx.Conn
}

func NewTestsRepository(conn *pgx.Conn) *TestsRepository {
	return &TestsRepository{conn: conn}
}

func (repo *TestsRepository) GetOne(testName string) (*models.Test, error) {
	sql := `select test_name, unit_type, "references", standards, available_modifiers from tests where test_name = $1`
	var test models.Test
	if err := repo.conn.QueryRow(context.Background(), sql, testName).Scan(&test.TestName, &test.UnitType, &test.References, &test.Standards, &test.AvailableModifiers); err != nil {
		return nil, err
	}
	return &test, nil
}

func (repo *TestsRepository) GetMany() (*[]models.Test, error) {
	sql := `select test_name, unit_type, "references", standards, available_modifiers from tests`
	rows, err := repo.conn.Query(context.Background(), sql)
	if err != nil {
		return nil, err
	}
	var tests []models.Test
	for rows.Next() {
		var test models.Test
		if err := rows.Scan(&test.TestName, &test.UnitType, &test.References, &test.Standards, &test.AvailableModifiers); err != nil {
			return nil, err
		}
		tests = append(tests, test)
	}
	return &tests, nil
}

func (repo *TestsRepository) Create(test *models.Test) error {
	sql := `
insert into tests (test_name, unit_type, "references", standards, available_modifiers)
values ($1, $2, $3, $4, $5)	
	`
	tag, err := repo.conn.Exec(context.Background(), sql, test.TestName, test.UnitType, test.References, test.Standards, test.AvailableModifiers)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("%v rows affected by insert (expected 1)", tag.RowsAffected())
	}
	return err
}

func (repo *TestsRepository) Update(test *models.Test) error {
	sql := `
update tests
set unit_type = $2, "references" = $3, standards = $4, available_modifiers = $5
where test_name = $1
	`
	tag, err := repo.conn.Exec(context.Background(), sql, test.TestName, test.UnitType, test.References, test.Standards, test.AvailableModifiers)
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("%v rows affected by update (expected 1)", tag.RowsAffected())
	}
	return err
}

func (repo *TestsRepository) Migrate() error {
	var currentVersion int
	row := repo.conn.QueryRow(context.Background(), "select current_version from table_versions where table_name = 'tests'")
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
create table tests (
	test_name text primary key,
	unit_type text not null,
	"references" text[] not null,
	standards text[] not null,
	available_modifiers text[] not null
)
		`)
		if err != nil {
			return err
		}
		_, err = repo.conn.Exec(context.Background(), `
insert into table_versions (table_name, current_version)
values ('tests', 1)
		`)
		if err != nil {
			return err
		}
	}
	return nil
}
