package dbstorage

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/metrics"
	"github.com/w1nsec/collector/internal/storage"
	"strconv"
	"strings"
)

type DBStorage interface {
	CheckConnection() error
	Close() error

	// check if tables created
	CreateTables() error

	storage.Storage
}

const (
	Counters = "counters"
	Gauges   = "gauges"
)

type postgresStorage struct {
	db  *sql.DB
	url string

	dbCtx context.Context

	addr     string
	username string
	password string
	dbName   string
}

func (pgStorage postgresStorage) Init() error {
	return pgStorage.CreateTables()
}

func (pgStorage postgresStorage) CheckStorage() error {
	return pgStorage.CheckConnection()
}

func (pgStorage postgresStorage) CreateTables() error {
	if pgStorage.db == nil {
		return fmt.Errorf("db not connected")
	}
	err := pgStorage.CheckConnection()
	if err != nil {
		return err
	}

	queryTb1 := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id varchar NOT NULL UNIQUE PRIMARY KEY,
		value int8 NULL
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
func (pgStorage postgresStorage) String() string {
	ms, err := pgStorage.GetAllMetrics()
	if err != nil {
		return ""
	}
	res := ""
	for _, m := range ms {
		var val float64
		switch m.MType {
		case metrics.Counter:
			val = float64(*m.Delta)
		case metrics.Gauge:
			val = *m.Value
		}
		mStr := fmt.Sprintf("Type: %s,\tID: %s,\tVal: %f<br>\n", m.MType, m.ID, val)
		res += mStr
	}

	return res
}

func (pgStorage postgresStorage) GetMetricString(mType, mName string) string {
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

func (pgStorage postgresStorage) GetMetric(mName string, mType string) (*metrics.Metrics, error) {
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
		return nil, nil
	}
	if row.Err() != nil {
		log.Error().Err(err).Send()
		return nil, nil
	}

	m := &metrics.Metrics{}
	switch mType {
	case metrics.Gauge:
		val, err := strconv.ParseFloat(*result, 64)
		if err != nil {
			fmt.Println(err)
			return nil, nil
		}
		m.ID = mName
		m.MType = metrics.Gauge
		m.Value = &val
	case metrics.Counter:
		val, err := strconv.ParseInt(*result, 10, 64)
		if err != nil {
			fmt.Println(err)
			return nil, nil
		}
		m.ID = mName
		m.MType = metrics.Counter
		m.Delta = &val
	}
	return m, nil

}

func (pgStorage postgresStorage) GetAllMetrics() ([]*metrics.Metrics, error) {
	// Gauges
	query := fmt.Sprintf("SELECT id, value from %s", Gauges)
	rows, err := pgStorage.db.QueryContext(pgStorage.dbCtx, query)
	if err != nil {
		return nil, nil
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
		return nil, nil
	}

	// Counters
	query = fmt.Sprintf("SELECT id, value from %s", Counters)
	rows, err = pgStorage.db.QueryContext(pgStorage.dbCtx, query)
	if err != nil {
		return nil, nil
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
		return nil, rows.Err()
	}

	return ms, nil
}

func NewStorage(url string) *postgresStorage {
	if !strings.Contains(url, "postgres://") {
		url = "postgres://" + url
	}
	db, err := sql.Open("pgx", url)
	if err != nil {
		log.Error().Err(err).Send()
		return nil
	}
	return &postgresStorage{
		// TODO set context to DB Storage as child context from Service
		dbCtx: context.Background(),
		db:    db,
		url:   url,
	}
}

func (pgStorage postgresStorage) Close() error {
	return pgStorage.db.Close()
}

func (pgStorage postgresStorage) CheckConnection() error {
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
