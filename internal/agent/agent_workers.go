package agent

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/metrics"
	"sync"
	"time"
)

func (agent Agent) generator(ctx context.Context,
	metricsChannel chan []*metrics.Metrics) {

	pollTicker := time.NewTicker(agent.pollInterval)
	for {
		select {
		case t1 := <-pollTicker.C:
			log.Info().
				Str("time", t1.Format(time.TimeOnly)).
				Msgf("Receiving started")

			localWG := &sync.WaitGroup{}
			localWG.Add(2)
			go func() {
				defer localWG.Done()
				agent.CollectMetrics(ctx)
			}()
			go func() {
				defer localWG.Done()
				agent.CollectGopsutilMetrics(ctx)
			}()

			// wait until all metrics not gathered
			localWG.Wait()
			log.Info().
				Str("time", t1.Format(time.TimeOnly)).
				Msgf("Receiving done")

			// get metrics from storage
			allMetrics, err := agent.store.GetAllMetrics(ctx)
			if err != nil {
				log.Error().Err(err).Send()
				continue
			}

			// add gathered metrics to channel
			go func() {
				metricsChannel <- allMetrics
			}()

		case <-ctx.Done():
			// writer should close channel
			close(metricsChannel)
			return
		}
	}
}

func (agent Agent) limiter(ctx context.Context,
	metricsChannel chan []*metrics.Metrics) {
	
	reportTicker := time.NewTicker(agent.reportInterval)
	var m sync.Mutex
	cond := sync.NewCond(&m)

	for i := 0; i < agent.rateLimit; i++ {
		go agent.worker(i, metricsChannel, cond)
	}

	// create workers
	for {
		select {
		case t2 := <-reportTicker.C:
			fmt.Println("Sending:", t2.Format(time.TimeOnly))
			fmt.Printf("Len: %d\n", len(metricsChannel))
			cond.Broadcast()

		case <-ctx.Done():
			return
		}
	}
}

func (agent Agent) worker(id int, jobs <-chan []*metrics.Metrics, c *sync.Cond) {
	// lock current worker
	c.L.Lock()
	for {
		// worker must sleep until report time
		// each worker should send ONLY ONE request to server
		c.Wait()
		// job == one metric batch
		job, ok := <-jobs
		//_, ok := <-jobs
		if !ok {
			// close worker, if jobs channel already closed
			return
		}

		// worker work
		log.Info().
			Int("worker", id).
			Msg("Sending")

		err := agent.SendBatch(job)
		if err != nil {
			log.Error().
				Int("worker", id).
				Err(err).Send()
			continue
		}
		log.Info().
			Int("worker", id).
			Msg("Done")

	}
}
