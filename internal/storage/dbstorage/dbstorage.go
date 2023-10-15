package dbstorage

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
	"strings"
)

type DBStorage interface {
	CheckConnection() error
	Close() error
}

type PostgresStorage struct {
	db  *sql.DB
	url string

	addr     string
	username string
	password string
	dbName   string
}

func NewPostgresStorage(url string) *PostgresStorage {
	if !strings.Contains(url, "postgres://") {
		url = "postgres://" + url
	}
	return &PostgresStorage{
		db:  nil,
		url: url,
	}
}

func (pgStorage PostgresStorage) Close() error {
	return pgStorage.db.Close()
}

func (pgStorage PostgresStorage) CheckConnection() error {
	//connectString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
	//	pgStorage.addr, pgStorage.username, pgStorage.password, pgStorage.dbName)

	var err error
	log.Info().Str("db_url", pgStorage.url).Send()
	pgStorage.db, err = sql.Open("pgx", pgStorage.url)
	if err != nil {
		log.Error().Err(err).Send()
		return err
	}
	err = pgStorage.db.Ping()
	if err != nil {
		return err
	}
	return nil
}
