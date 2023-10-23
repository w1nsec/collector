package dbstorage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/metrics"
)

func (pgStorage postgresStorage) UpdateCounters(name string, value int64) error {
	query := fmt.Sprintf("update %s set value = value + $1 where id = $2", Counters)
	result, err := pgStorage.db.ExecContext(pgStorage.dbCtx, query, value, name)
	if err != nil {
		return err
	}
	if result != nil {
		num, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if num != 1 {
			return fmt.Errorf("expected to affect 1 row, affected %d", num)
		}
	}
	return nil
}

func (pgStorage postgresStorage) UpdateGauges(name string, value float64) error {
	query := fmt.Sprintf("update %s set value = $1 where id = $2", Gauges)
	result, err := pgStorage.db.ExecContext(pgStorage.dbCtx, query, value, name)
	if err != nil {
		return err
	}
	if result != nil {
		num, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if num != 1 {
			return fmt.Errorf("expected to affect 1 row, affected %d", num)
		}
	}

	return nil
}

func (pgStorage postgresStorage) UpdateMetric(newMetric *metrics.Metrics) error {
	var result sql.Result
	var err error
	switch newMetric.MType {
	case metrics.Gauge:
		query := fmt.Sprintf("update %s set value = $1 where id = $2", Gauges)
		result, err = pgStorage.db.ExecContext(pgStorage.dbCtx, query, newMetric.Value, newMetric.ID)
	case metrics.Counter:
		query := fmt.Sprintf("update %s set value = value + $1 where id = $2", Counters)
		result, err = pgStorage.db.ExecContext(pgStorage.dbCtx, query, newMetric.Delta, newMetric.ID)
	}
	if err != nil {
		return err
	}

	if result != nil {
		num, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if num != 1 {
			//return fmt.Errorf("expected to affect 1 row, affected %d", num)

			// try to add metric
			err = pgStorage.AddMetric(newMetric)
			return err
		}
	}
	return nil
}

func (pgStorage postgresStorage) UpdateMetrics(ctx context.Context, newMetrics []*metrics.Metrics) error {
	var (
		queryGauge   = fmt.Sprintf("insert into %s(id, value) values ($1, $2) on conflict (id) do update set value = $2", Gauges)
		queryCounter = fmt.Sprintf("insert into %s(id, value) values ($1, $2) on conflict (id) do update set value = %s.value + $2", Counters, Counters)
	)

	// start transaction
	tx, err := pgStorage.db.Begin()
	if err != nil {
		return err
	}

	// insert metrics to DB
	for _, newMetric := range newMetrics {
		switch newMetric.MType {
		case metrics.Gauge:
			_, err = tx.ExecContext(ctx, queryGauge, newMetric.ID, newMetric.Value)
		case metrics.Counter:
			_, err = tx.ExecContext(ctx, queryCounter, newMetric.ID, newMetric.Delta)
		}
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// commit changes to DB
	return tx.Commit()
}

func (pgStorage postgresStorage) AddMetric(newMetric *metrics.Metrics) error {
	var result sql.Result
	var err error
	var query = "INSERT into %s (id, value) values ($1, $2)"
	switch newMetric.MType {
	case metrics.Gauge:
		query = fmt.Sprintf(query, Gauges)
		result, err = pgStorage.db.ExecContext(pgStorage.dbCtx, query, newMetric.ID, newMetric.Value)
	case metrics.Counter:
		query = fmt.Sprintf(query, Counters)
		result, err = pgStorage.db.ExecContext(pgStorage.dbCtx, query, newMetric.ID, newMetric.Delta)
	}
	log.Info().Msgf(query)
	if err != nil {
		return err
	}
	if result != nil {
		num, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if num != 1 {
			return fmt.Errorf("expected to affect 1 row, affected %d", num)
		}
	}
	return nil
}
