package agent

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/metrics"
	"sync"
	"syscall"
	"time"
)

func (agent Agent) generator(ctx context.Context,
	metricsChannel chan []*metrics.Metrics) {
	// writer should close channel
	defer close(metricsChannel)

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
			log.Info().
				Str("func", "generator").
				Msg("Closing")
			return
		}
	}
}

func (agent Agent) limiter(ctx context.Context,
	metricsChannel chan []*metrics.Metrics) {

	reportTimer := time.NewTimer(agent.reportInterval)

	// waiting for report time
	for t2 := range reportTimer.C {
		fmt.Println("Sending:", t2.Format(time.TimeOnly))
		// create workers
		for i := 0; i < agent.rateLimit; i++ {
			go agent.worker(i, metricsChannel)
		}
	}

	<-ctx.Done()
	log.Info().
		Str("func", "limiter").
		Msg("Closing")

}

func (agent Agent) worker(id int, jobs <-chan []*metrics.Metrics) {
	for {
		select {
		case job := <-jobs:
			// worker work
			log.Info().
				Int("worker", id).
				Msg("Sending")

			err := agent.SendBatch(job)
			if err != nil {
				log.Error().
					Int("worker", id).
					Err(err).Send()
				agent.errorsCh <- err
				continue
			}
			log.Info().
				Int("worker", id).
				Msg("Done")
			agent.successReport <- struct{}{}

		// if localerrors too many, go to sleep,
		case <-agent.sleepCh[id]:
			time.Sleep(time.Second * time.Duration(10*agent.sleepCount))
		}

	}
}

func (agent Agent) validateErrors(ctx context.Context) {
	var (
		curErrCount = 0
		maxErrCount = int(agent.retryCount) * agent.rateLimit
	)

	for {
		select {
		// check connection
		case err := <-agent.errorsCh:
			if errors.Is(err, syscall.ECONNREFUSED) {
				curErrCount++
			}
			log.Info().Msgf("Errors count: %d/%d", curErrCount, maxErrCount)
			// localerrors should be more, as frequency sending increase
			if curErrCount == maxErrCount {
				/*
					increase sleepCount (for sleep time in workers)
					but not to infinity
				*/
				if agent.sleepCount < agent.retryCount {
					agent.sleepCount++
				}

				// send sleep signal to all workers
				wg := sync.WaitGroup{}
				wg.Add(1)
				// deadlock bypass
				go func() {
					defer wg.Done()
					for _, ch := range agent.sleepCh {
						ch <- struct{}{}
					}
				}()
				wg.Wait()
			}
		case <-agent.successReport:
			// reset count of localerrors
			curErrCount = 0

			// reset count of sleep times
			agent.sleepCount = 0
		case <-ctx.Done():
			log.Info().
				Str("func", "error-validator").
				Msg("Closing")
			return
		}
	}
}
