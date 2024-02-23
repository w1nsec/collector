package agent

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/w1nsec/collector/internal/config"
)

const (
	defaultPollInterval   = 2 * time.Second
	defaultReportInterval = 10 * time.Second
)

// JSON config example
//{
//	"address": "localhost:8080", 		// аналог переменной окружения ADDRESS или флага -a
//	"report_interval": "1s",     		// аналог переменной окружения REPORT_INTERVAL или флага -r
//	"poll_interval": "1s",           	// аналог переменной окружения POLL_INTERVAL или флага -p
//	"crypto_key": "/path/to/key.pem" 	// аналог переменной окружения CRYPTO_KEY или флага -crypto-key
//}

type Args struct {
	Addr           string `json:"address"`
	PollInterval   int    `json:"poll_interval"`
	ReportInterval int    `json:"report_interval"`
	Key            string
	Rate           int
	CryptoKey      string `json:"crypto_key"`

	LogLevel string

	Protocol string `json:"protocol"`
}

// ReadConfig fill Args struct (for client)
func ReadConfig(path string) (conf *Args, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	conf = new(Args)
	err = json.NewDecoder(file).Decode(&conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

// AgentSelectArgs return params in Args struct, that need for agent successfully run
func AgentSelectArgs() *Args {
	// check flags
	var (
		args                        *Args
		flagAddr, flagKey           string
		flagCryptoKey, flagConfig   string
		flagPoll, flagRep, flagRate int
		flagProto                   string
	)
	flag.StringVar(&flagAddr, "a", "localhost:8080",
		"address for metric transport")
	flag.IntVar(&flagPoll, "p", int(defaultPollInterval.Seconds()),
		"frequency of gathering metrics")
	flag.IntVar(&flagRep, "r", int(defaultReportInterval.Seconds()),
		"frequency of sending metrics")
	flag.StringVar(&flagKey, "k", "",
		"salt for hmac")
	flag.IntVar(&flagRate, "l", 2,
		"max goroutines count")
	flag.StringVar(&flagCryptoKey, "crypto-key", "",
		"rsa public key path (in pem format), used for encrypt messages")
	flag.StringVar(&flagConfig, "config", "",
		"path to config file")
	flag.StringVar(&flagProto, "proto", "http",
		"default")
	flag.Parse()

	// Read config file
	// priority to ENV variable
	cpath := os.Getenv(config.ConfFile)
	if cpath != "" {
		flagConfig = cpath
	}

	if flagConfig != "" {
		var err error
		args, err = ReadConfig(flagConfig)
		if err != nil {
			log.Debug().Err(err).Send()
		}
	}

	// If we setup Args early, then they will rewrite later
	// check that args have been filled
	if args == nil {
		args = new(Args)
	}

	// set log level
	args.LogLevel = config.DefaultLogLevel

	// get OS Environment variables
	if flagAddr != "" {
		args.Addr = flagAddr
	}
	// priority to ENV variables
	addr := os.Getenv(config.Address)
	if addr != "" {
		args.Addr = addr
	}

	if flagPoll != 0 {
		args.PollInterval = flagPoll
	}
	var err error
	pollInterval, err := strconv.Atoi(os.Getenv(config.PollInterval))
	if err == nil {
		args.PollInterval = pollInterval
	}

	if flagRep != 0 {
		args.ReportInterval = flagRep
	}
	reportInterval, err := strconv.Atoi(os.Getenv(config.ReportInterval))
	if err == nil {
		args.ReportInterval = reportInterval
	}

	if flagKey != "" {
		args.Key = flagKey
	}
	key := os.Getenv(config.Key)
	if key != "" {
		args.Key = key
	}

	// increment15
	if flagRate != 0 {
		args.Rate = flagRate
	}
	rate, err := strconv.Atoi(os.Getenv(config.RateLimit))
	if err == nil {
		args.Rate = rate
	}

	// increment 21
	if flagCryptoKey != "" {
		args.CryptoKey = flagCryptoKey
	}
	cryptoKey := os.Getenv(config.CryptoKey)
	if cryptoKey != "" {
		args.CryptoKey = cryptoKey
	}

	// increment 25
	if flagProto != "" {
		args.Protocol = flagProto
	}
	proto := os.Getenv(config.Protocol)
	if proto != "" {
		args.Protocol = proto
	}
	if args.Protocol != config.ProtoGRPC && args.Protocol != config.ProtoHTTP {
		args.Protocol = config.ProtoHTTP
	}

	return args
}
