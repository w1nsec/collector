// Package agent contains Agent struct
// that provide functionality for collecting and sending
// metrics to server part
package agent

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	config "github.com/w1nsec/collector/internal/config/agent"
	"github.com/w1nsec/collector/internal/metrics"
	"github.com/w1nsec/collector/internal/storage"
	"github.com/w1nsec/collector/internal/storage/memstorage"
	"github.com/w1nsec/collector/internal/utils/ip"
	"github.com/w1nsec/go-examples/crypto"
)

// Params for sleeping agent if it receives too many errors
var (
	timeout = 10 * time.Second
	//retryStep     = uint(2) // 2 seconds
	maxRetryCount = uint(5)
	buildVersion  = "N/A"
	buildDate     = "N/A"
	buildCommit   = "N/A"
)

// Agent struct, that contains Storage and other config options for running agent
type Agent struct {
	addr         net.Addr
	metricsPoint string

	store          storage.Storage
	pollInterval   time.Duration
	reportInterval time.Duration
	compression    bool
	httpClient     *http.Client

	// increment 13
	retryCount uint
	sleepCount uint
	sleepCh    []chan struct{}

	// param needs for signing send body (increment 14)
	secret string

	// increment 15
	rateLimit     int
	errorsCh      chan error
	successReport chan struct{}

	// increment 21
	pubKey *rsa.PublicKey

	// increment 24
	realIP []string
}

// NewAgent is constructor for Agent struct
func NewAgent(args *config.Args) (*Agent, error) {
	netAddr, err := net.ResolveTCPAddr("tcp", args.Addr)
	if err != nil {
		return nil, err
	}

	sleepCh := make([]chan struct{}, args.Rate)
	for id := range sleepCh {
		sleepCh[id] = make(chan struct{})
	}

	agent := &Agent{
		addr:           netAddr,
		metricsPoint:   "update",
		pollInterval:   time.Duration(args.PollInterval) * time.Second,
		reportInterval: time.Duration(args.ReportInterval) * time.Second,

		store:      memstorage.NewMemStorage(),
		httpClient: &http.Client{Timeout: timeout},

		// increment 13
		retryCount: maxRetryCount,
		sleepCount: 0,
		sleepCh:    sleepCh,

		// increment 14
		secret: args.Key,

		// increment 15
		rateLimit:     args.Rate,
		errorsCh:      make(chan error, args.Rate),
		successReport: make(chan struct{}, args.Rate),
	}

	// increment 21
	if args.CryptoKey != "" {
		key, err := os.ReadFile(args.CryptoKey)
		if err != nil {
			return nil, err
		}

		agent.pubKey, err = crypto.ReadPubKey(key)
		if err != nil {
			return nil, err
		}
	}

	agent.compression = true
	return agent, nil
}

// Start run Agent functionality (collect and send metrics)
func (agent Agent) Start(ctx context.Context) error {
	fmt.Printf(
		"Build version: %s\n"+
			"Build daate: %s\n"+
			"Build commit: %s\n",
		buildVersion,
		buildDate,
		buildCommit,
	)

	// increment24 gel local IP
	ips, err := ip.GetIPv4()
	if err != nil {
		return err
	}
	agent.realIP = ips

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	metricsChannel := make(chan []*metrics.Metrics)

	// fill agent store
	go func() {
		agent.generator(ctx, metricsChannel)
	}()

	// create workers pull
	go func() {
		agent.limiter(ctx, metricsChannel)
	}()

	// validate localerrors count
	go func() {
		agent.validateErrors(ctx)
	}()

	// waiting, until goroutines not done
	<-ctx.Done()

	return nil
}

// Close channels
func (agent Agent) Close() {
	close(agent.errorsCh)
	close(agent.successReport)
	err := agent.store.Close(context.TODO())
	if err != nil {
		log.Error().Err(err).Send()
	}
}
