package repositories

import (
	"config/models"
	"context"

	"github.com/jackc/pgx/v5"
)

type UnitsRepository struct {
	conn *pgx.Conn
}

func NewUnitsRepository(conn *pgx.Conn) *UnitsRepository {
	return &UnitsRepository{conn: conn}
}

func (repo *UnitsRepository) GetMany() (*[]models.Unit, error) {
	sql := "select full_name, full_name_plural, abbreviation, measurement_system, unit_type from units"
	rows, err := repo.conn.Query(context.Background(), sql)
	if err != nil {
		return nil, err
	}
	var units []models.Unit
	for rows.Next() {
		var unit models.Unit
		if err := rows.Scan(&unit.FullName, &unit.FullNamePlural, &unit.Abbreviation, &unit.MeasurementSystem, &unit.UnitType); err != nil {
			return nil, err
		}
		units = append(units, unit)
	}
	return &units, nil
}

func (repo *UnitsRepository) InsertMany(units *[]models.Unit, by string, tx pgx.Tx) error {
	sql := `
insert into units (full_name, full_name_plural, abbreviation, measurement_system, unit_type)	
values ($1, $2, $3, $4, $5)
	`
	for _, unit := range *units {
		if tx != nil {
			if _, err := tx.Exec(context.Background(), sql, unit.FullName, unit.FullNamePlural, unit.Abbreviation, unit.MeasurementSystem, unit.UnitType); err != nil {
				return err
			}
		} else {
			if _, err := repo.conn.Exec(context.Background(), sql, unit.FullName, unit.FullNamePlural, unit.Abbreviation, unit.MeasurementSystem, unit.UnitType); err != nil {
				return err
			}
		}
	}
	return nil
}

func (repo *UnitsRepository) Migrate() error {
	var currentVersion int
	row := repo.conn.QueryRow(context.Background(), "select current_version from table_versions where table_name = 'units'")
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
create table units (
	full_name text primary key,
	full_name_plural text not null,
	abbreviation text not null,
	measurement_system text not null,
	unit_type text not null
)
		`)
		if err != nil {
			return err
		}
		_, err = repo.conn.Exec(context.Background(), `
insert into table_versions (table_name, current_version)
values ('units', 1)
		`)
		if err != nil {
			return err
		}
	}
	return nil
}
