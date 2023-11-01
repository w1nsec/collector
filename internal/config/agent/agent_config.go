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

	LogLevel string
}

func AgentSelectArgs(args *Args) {
	// set log level
	args.LogLevel = "info"

	// get OS Environment variables
	args.Addr = os.Getenv("ADDRESS")
	envPoll, err := strconv.Atoi(os.Getenv("POLL_INTERVAL"))
	if err == nil {
		args.PollInterval = envPoll
	}
	envRep, err := strconv.Atoi(os.Getenv("REPORT_INTERVAL"))
	if err == nil {
		args.ReportInterval = envRep
	}

	args.Key = os.Getenv("KEY")

	// check flags
	var (
		flagAddr, flagKey string
		flagPoll, flagRep int
	)
	flag.StringVar(&flagAddr, "a", "localhost:8080",
		"address for metric server")
	flag.IntVar(&flagPoll, "p", int(defaultPollInterval.Seconds()),
		"frequency of gathering metrics")
	flag.IntVar(&flagRep, "r", int(defaultReportInterval.Seconds()),
		"frequency of sending metrics")
	flag.StringVar(&flagKey, "k", "",
		"salt for hmac")

	flag.Parse()

	if args.Addr == "" {
		args.Addr = flagAddr
	}

	if args.PollInterval == 0 {
		args.PollInterval = flagPoll
	}
	if args.ReportInterval == 0 {
		args.ReportInterval = flagRep
	}

	if args.Key == "" {
		args.Key = flagKey
	}

}
