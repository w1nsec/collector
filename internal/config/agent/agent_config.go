package agent

import (
	"flag"
	"os"
	"strconv"
	"time"
)

const (
	defaultPollInterval   = 2 * time.Second
	defaultReportInterval = 10 * time.Second
)

type Args struct {
	Addr           string
	PollInterval   int
	ReportInterval int
	Key            string
	Rate           int

	LogLevel string
}

func AgentSelectArgs(args *Args) {
	// set log level
	args.LogLevel = "info"

	// check flags
	var (
		flagAddr, flagKey           string
		flagPoll, flagRep, flagRate int
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

	flag.Parse()

	// get OS Environment variables
	args.Addr = os.Getenv("ADDRESS")
	if args.Addr == "" {
		args.Addr = flagAddr
	}

	var err error
	args.PollInterval, err = strconv.Atoi(os.Getenv("POLL_INTERVAL"))
	if err != nil {
		args.PollInterval = flagPoll
	}

	args.ReportInterval, err = strconv.Atoi(os.Getenv("REPORT_INTERVAL"))
	if err != nil {
		args.ReportInterval = flagRep
	}

	args.Key = os.Getenv("KEY")
	if args.Key == "" {
		args.Key = flagKey
	}

	// increment15
	args.Rate, err = strconv.Atoi(os.Getenv("RATE_LIMIT"))
	if err != nil {
		args.Rate = flagRate
	}

}
