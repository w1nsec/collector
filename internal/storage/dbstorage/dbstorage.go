package dbstorage

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/metrics"
	"github.com/w1nsec/collector/internal/storage/memstorage"
	"strconv"
	"strings"
)

type DBStorage interface {
	CheckConnection() error
	Close() error

	// check if tables created
	CreateTables() error

	memstorage.Storage
}

const (
	Counters = "counters"
	Gauges   = "gauges"
)

type PostgresStorage struct {
	db  *sql.DB
	url string

	dbCtx context.Context

	addr     string
	username string
	password string
	dbName   string
}

func (pgStorage PostgresStorage) CreateTables() error {
	if pgStorage.db == nil {
		return fmt.Errorf("db not connected")
	}
	err := pgStorage.CheckConnection()
	if err != nil {
		return err
	}

	queryTb1 := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id varchar NOT NULL UNIQUE PRIMARY KEY,
		value integer NULL
	);`, Counters)

	queryTb2 := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		id varchar NOT NULL UNIQUE PRIMARY KEY ,
		value double precision NULL
	);`, Gauges)

	_, err = pgStorage.db.ExecContext(pgStorage.dbCtx, queryTb1)
	if err != nil {
		return err
	}

	_, err = pgStorage.db.ExecContext(pgStorage.dbCtx, queryTb2)
	if err != nil {
		return err
	}

	return nil
}

// interface Storage
func (pgStorage PostgresStorage) UpdateCounters(name string, value int64) {
	query := fmt.Sprintf("update %s set value = $1 where id = $2", Counters)
	result, err := pgStorage.db.ExecContext(pgStorage.dbCtx, query, value, name)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
	if result != nil {
		num, err := result.RowsAffected()
		if err != nil {
			log.Error().Err(err).Send()
			return
		}
		if num != 1 {
			log.Error().
				Err(fmt.Errorf("expected to affect 1 row, affected %d", num)).
				Send()
			return
		}
	}
}

func (pgStorage PostgresStorage) UpdateGauges(name string, value float64) {
	query := fmt.Sprintf("update %s set value = $1 where id = $2", Gauges)
	result, err := pgStorage.db.ExecContext(pgStorage.dbCtx, query, value, name)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
	if result != nil {
		num, err := result.RowsAffected()
		if err != nil {
			log.Error().Err(err).Send()
			return
		}
		if num != 1 {
			log.Error().
				Err(fmt.Errorf("expected to affect 1 row, affected %d", num)).
				Send()
			return
		}
	}

}

func (pgStorage PostgresStorage) String() string {
	//TODO implement me
	panic("implement me")
}

func (pgStorage PostgresStorage) GetMetricString(mType, mName string) string {
	tbName := Counters
	if mType == metrics.Gauge {
		tbName = Gauges
	}
	query := fmt.Sprintf("SELECT DISTINCT value FROM %s where id = $1", tbName)
	row := pgStorage.db.QueryRowContext(pgStorage.dbCtx, query, mName)

	result := new(string)
	err := row.Scan(result)
	if err != nil {
		log.Error().Err(err).Send()
		return ""
	}
	if row.Err() != nil {
		log.Error().Err(err).Send()
		return ""
	}

	return *result

}

func (pgStorage PostgresStorage) GetMetric(mName string, mType string) *metrics.Metrics {
	tbName := Counters
	if mType == metrics.Gauge {
		tbName = Gauges
	}
	query := fmt.Sprintf("SELECT DISTINCT value FROM %s where id = $1", tbName)
	row := pgStorage.db.QueryRowContext(pgStorage.dbCtx, query, mName)

	result := new(string)
	err := row.Scan(result)
	if err != nil {
		log.Error().Err(err).Send()
		return nil
	}
	if row.Err() != nil {
		log.Error().Err(err).Send()
		return nil
	}

	m := &metrics.Metrics{}
	switch mType {
	case metrics.Gauge:
		val, err := strconv.ParseFloat(*result, 64)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		m.ID = mName
		m.Value = &val
	case metrics.Counter:
		val, err := strconv.ParseInt(*result, 10, 64)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		m.ID = mName
		m.Delta = &val
	}
	return m

}

func (pgStorage PostgresStorage) UpdateMetric(newMetric *metrics.Metrics) {
	var result sql.Result
	var err error
	switch newMetric.ID {
	case metrics.Gauge:
		query := fmt.Sprintf("update %s set value = $1 where id = $2", Gauges)
		result, err = pgStorage.db.ExecContext(pgStorage.dbCtx, query, newMetric.Value, newMetric.ID)
	case metrics.Counter:
		query := fmt.Sprintf("update %s set value = $1 where id = $2", Counters)
		result, err = pgStorage.db.ExecContext(pgStorage.dbCtx, query, newMetric.Delta, newMetric.ID)
	}
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	if result != nil {
		num, err := result.RowsAffected()
		if err != nil {
			log.Error().Err(err).Send()
			return
		}
		if num != 1 {
			log.Error().
				Err(fmt.Errorf("expected to affect 1 row, affected %d", num)).
				Send()
			return
		}
	}
}

func (pgStorage PostgresStorage) AddMetric(newMetric *metrics.Metrics) {
	var result sql.Result
	var err error
	switch newMetric.ID {
	case metrics.Gauge:
		query := fmt.Sprintf("INSERT %s (id, value) values ($1, $2)", Gauges)
		result, err = pgStorage.db.ExecContext(pgStorage.dbCtx, query, newMetric.ID, newMetric.Value)
	case metrics.Counter:
		query := fmt.Sprintf("INSERT %s (id, value) values ($1, $2)", Counters)
		result, err = pgStorage.db.ExecContext(pgStorage.dbCtx, query, newMetric.ID, newMetric.Delta)
	}
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
	if result != nil {
		num, err := result.RowsAffected()
		if err != nil {
			log.Error().Err(err).Send()
			return
		}
		if num != 1 {
			log.Error().
				Err(fmt.Errorf("expected to affect 1 row, affected %d", num)).
				Send()
			return
		}
	}
}

func (pgStorage PostgresStorage) GetAllMetrics() []*metrics.Metrics {
	// Gauges
	query := fmt.Sprintf("SELECT id, value from %s", Gauges)
	rows, err := pgStorage.db.QueryContext(pgStorage.dbCtx, query)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var ms = make([]*metrics.Metrics, 0)

	id := new(string)
	valGauge := new(float64)

	for rows.Next() {
		err := rows.Scan(id, valGauge)
		if err != nil {
			fmt.Println(err)
			continue
		}
		m := &metrics.Metrics{
			ID:    *id,
			Value: valGauge,
		}
		ms = append(ms, m)
	}

	if rows.Err() != nil {
		return nil
	}

	// Counters
	query = fmt.Sprintf("SELECT id, value from %s", Counters)
	rows, err = pgStorage.db.QueryContext(pgStorage.dbCtx, query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	valCounter := new(int64)
	for rows.Next() {
		err := rows.Scan(id, valCounter)
		if err != nil {
			fmt.Println(err)
			continue
		}
		m := &metrics.Metrics{
			ID:    *id,
			Delta: valCounter,
		}
		ms = append(ms, m)
	}

	if rows.Err() != nil {
		return nil
	}

	return ms
}

func NewPostgresStorage(url string) *PostgresStorage {
	if !strings.Contains(url, "postgres://") {
		url = "postgres://" + url
	}
	return &PostgresStorage{
		// TODO set context to DB Storage as child context from Service
		dbCtx: context.Background(),
		db:    nil,
		url:   url,
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
