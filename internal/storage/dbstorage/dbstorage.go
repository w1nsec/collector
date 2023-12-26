package dbstorage

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/metrics"
	"github.com/w1nsec/collector/internal/storage"
)

type DBStorage interface {
	CheckConnection() error
	Close(context.Context) error

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

	_, err = pgStorage.db.Exec(queryTb1)
	if err != nil {
		return err
	}

	_, err = pgStorage.db.Exec(queryTb2)
	if err != nil {
		return err
	}

	return nil
}

// interface Storage
func (pgStorage postgresStorage) String(ctx context.Context) string {
	ms, err := pgStorage.GetAllMetrics(ctx)
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

func (pgStorage postgresStorage) GetMetricString(ctx context.Context, mType, mName string) string {
	tbName := Counters
	if mType == metrics.Gauge {
		tbName = Gauges
	}
	query := fmt.Sprintf("SELECT DISTINCT value FROM %s where id = $1", tbName)
	row := pgStorage.db.QueryRowContext(ctx, query, mName)

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

func (pgStorage postgresStorage) GetMetric(ctx context.Context, mName string, mType string) (*metrics.Metrics, error) {
	tbName := Counters
	if mType == metrics.Gauge {
		tbName = Gauges
	}
	query := fmt.Sprintf("SELECT DISTINCT value FROM %s where id = $1", tbName)
	row := pgStorage.db.QueryRowContext(ctx, query, mName)

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
			log.Error().Err(err).Send()
			return nil, nil
		}
		m.ID = mName
		m.MType = metrics.Gauge
		m.Value = &val
	case metrics.Counter:
		val, err := strconv.ParseInt(*result, 10, 64)
		if err != nil {
			log.Error().Err(err).Send()
			return nil, nil
		}
		m.ID = mName
		m.MType = metrics.Counter
		m.Delta = &val
	}
	return m, nil

}

func (pgStorage postgresStorage) GetAllMetrics(ctx context.Context) ([]*metrics.Metrics, error) {
	var (
		ms = make([]*metrics.Metrics, 0)
	)

	msCounters, err := pgStorage.getCounters(ctx)
	if err != nil {
		log.Error().Err(err).Send()
	}

	msGauges, err := pgStorage.getGauges(ctx)
	if err != nil {
		log.Error().Err(err).Send()
	}

	ms = append(msCounters, msGauges...)
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
		db:  db,
		url: url,
	}
}

func (pgStorage postgresStorage) Close(context.Context) error {
	return pgStorage.db.Close()
}

func (pgStorage postgresStorage) CheckConnection() error {
	//connectString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
	//	pgStorage.addr, pgStorage.username, pgStorage.password, pgStorage.dbName)

	var err error
	log.Info().Str("db_url", pgStorage.url).Send()
	if pgStorage.db == nil {
		pgStorage.db, err = sql.Open("pgx", pgStorage.url)
		if err != nil {
			log.Error().Err(err).Send()
			return err
		}
	}

	err = pgStorage.db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (pgStorage postgresStorage) getGauges(ctx context.Context) ([]*metrics.Metrics, error) {
	// Gauges
	query := fmt.Sprintf("SELECT id, value from %s", Gauges)
	rows, err := pgStorage.db.QueryContext(ctx, query)
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	var (
		ms       = make([]*metrics.Metrics, 0)
		id       = new(string)
		valGauge = new(float64)
	)

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
		return nil, rows.Err()
	}
	return ms, nil
}

func (pgStorage postgresStorage) getCounters(ctx context.Context) ([]*metrics.Metrics, error) {
	// Counters
	query := fmt.Sprintf("SELECT id, value from %s", Counters)
	rows, err := pgStorage.db.QueryContext(ctx, query)
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	var (
		ms         = make([]*metrics.Metrics, 0)
		id         = new(string)
		valCounter = new(int64)
	)

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
