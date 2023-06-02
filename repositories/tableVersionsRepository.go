package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type TableVersionsRepository struct {
	conn *pgx.Conn
}

func NewTableVersionsRepository(conn *pgx.Conn) *TableVersionsRepository {
	return &TableVersionsRepository{conn: conn}
}

func (repo *TableVersionsRepository) Migrate() error {
	var exists bool
	err := repo.conn.QueryRow(context.Background(), "select exists (select 1 from pg_tables where tablename = 'table_versions')").Scan(&exists)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Info().Msg("table_versions not found")
			exists = false
		} else {
			return err
		}
	}
	var currentVersion int
	if exists {
		log.Info().Msg("table_versions found")
		err = repo.conn.QueryRow(context.Background(), "select current_version from table_versions where table_name = 'table_versions'").Scan(&currentVersion)
		if err != nil {
			if err == pgx.ErrNoRows {
				currentVersion = 0
			} else {
				return err
			}
		}
	} else {
		currentVersion = 0
	}
	if currentVersion < 1 {
		log.Info().Msg("creating table table_versions")
		_, err = repo.conn.Exec(context.Background(), `
create table table_versions (
	table_name text primary key,
	current_version int not null
)
		`)
		if err != nil {
			return err
		}
		_, err = repo.conn.Exec(context.Background(), `
insert into table_versions (table_name, current_version)
values ('table_versions', 1)
		`)
		if err != nil {
			return err
		}
	}
	return nil
}
