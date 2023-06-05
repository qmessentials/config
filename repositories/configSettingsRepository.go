package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type ConfigSettingsRepository struct {
	conn *pgx.Conn
}

func NewConfigSettingsRepository(conn *pgx.Conn) *ConfigSettingsRepository {
	return &ConfigSettingsRepository{conn: conn}
}

func (repo *ConfigSettingsRepository) GetUnderlyingConnection() *pgx.Conn {
	return repo.conn
}

func (repo *ConfigSettingsRepository) GetOne(name string) (*[]string, error) {
	var values []string
	if err := repo.conn.QueryRow(
		context.Background(),
		"select setting_values from config_settings where name = $1", name).Scan(&values); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &values, nil
}

func (repo *ConfigSettingsRepository) GetOneFlag(name string) (bool, error) {
	var flag bool
	if err := repo.conn.QueryRow(
		context.Background(),
		"select setting_values[1]::boolean as flag from config_settings where name = $1", name).Scan(&flag); err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return flag, nil
}

func (repo *ConfigSettingsRepository) Upsert(name string, values []string, tx pgx.Tx) error {
	sql := `
insert into config_settings (name, setting_values) 
	values ($1, $2)
on conflict (name) do update
set setting_values = $2
	`
	if tx != nil {
		_, err := tx.Exec(context.Background(), sql, name, values)
		return err
	} else {
		_, err := repo.conn.Exec(context.Background(), sql, name, values)
		return err
	}
}

func (repo *ConfigSettingsRepository) Migrate() error {
	var currentVersion int
	row := repo.conn.QueryRow(context.Background(), "select current_version from table_versions where table_name = 'config_settings'")
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
create table config_settings (
	name text primary key,
	setting_values text[] not null
)
		`)
		if err != nil {
			return err
		}
		_, err = repo.conn.Exec(context.Background(), `
insert into table_versions (table_name, current_version)
values ('config_settings', 1)		
		`)
		if err != nil {
			return err
		}
	}
	return nil
}
