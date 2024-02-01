package server

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

// JSON config example
//{
//	"address": "localhost:8080", 		// аналог переменной окружения ADDRESS или флага -a
//	"restore": true, 					// аналог переменной окружения RESTORE или флага -r
//	"store_interval": "1s", 			// аналог переменной окружения STORE_INTERVAL или флага -i
//	"store_file": "/path/to/file.db", 	// аналог переменной окружения STORE_FILE или -f
//	"database_dsn": "", 				// аналог переменной окружения DATABASE_DSN или флага -d
//	"crypto_key": "/path/to/key.pem" 	// аналог переменной окружения CRYPTO_KEY или флага -crypto-key
//}

type Args struct {
	Addr     string `json:"address"`
	LogLevel string

	// increment 9, FILE_STORAGE
	StoreInterval uint64 `json:"store_interval"`
	StoragePath   string `json:"store_file"`
	Restore       bool   `json:"restore"`

	// increment 10 DB
	DatabaseURL string `json:"database_dsn"`

	// increment 14
	Key string

	// increment 21
	CryptoKey string `json:"crypto_key"`
}

// ReadConfig fill Args struct (for server)
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

// ServerArgsParse return params in Args struct, that need for server successfully run
func ServerArgsParse() *Args {

	// read flags
	var (
		args         *Args
		flagAddr     string
		flagLogLevel string

		// increment 9
		flagStoreInterval uint64
		flagStoragePath   string
		flagRestore       bool

		// increment 10
		flagDatabaseStr string

		// increment 14
		flagKey string

		// increment 21
		flagCryptoKey string
		flagConfig    string
	)

	flag.StringVar(&flagAddr, "a", "localhost:8080",
		"address for transport")
	flag.StringVar(&flagLogLevel, "l", "info",
		"log level")

	// increment 9, FILE_STORAGE
	flag.Uint64Var(&flagStoreInterval, "i", 300,
		"interval in seconds for write store data to file")
	flag.StringVar(&flagStoragePath, "f", "/tmp/metrics-db.json",
		"file for saving metrics")
	flag.BoolVar(&flagRestore, "r", true,
		"restore from file-db on startup")

	// increment 10, connect to DB
	flag.StringVar(&flagDatabaseStr, "d", "", "DB connect string")

	// increment 14, generate hash for requests body
	flag.StringVar(&flagKey, "k", "", "salt for hmac")

	// increment 21, decrypt requests
	flag.StringVar(&flagCryptoKey, "crypto-key", "",
		"rsa private key path (in pem format), used for encrypt messages")
	flag.StringVar(&flagConfig, "config", "",
		"path to config file")

	flag.Parse()

	// Read config file
	// priority to ENV variable
	cpath := os.Getenv("CONFIG")
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

	addr := os.Getenv("ADDRESS")
	if addr != "" {
		args.Addr = addr
	}
	args.LogLevel = os.Getenv("LOGLEVEL")

	// increment 9
	storeInterval, err := strconv.ParseUint(
		os.Getenv("STORE_INTERVAL"), 10, 64)
	if err == nil {
		args.StoreInterval = storeInterval
	}

	storagePath := os.Getenv("FILE_STORAGE_PATH")
	if storagePath != "" {
		args.StoragePath = storagePath
	}

	restore, err := strconv.ParseBool(os.Getenv("RESTORE"))
	if err == nil {
		args.Restore = restore
	}

	// increment 10
	databaseURL := os.Getenv("DATABASE_DSN")
	if databaseURL != "" {
		args.DatabaseURL = databaseURL
	}

	// increment 14
	key := os.Getenv("KEY")
	if key != "" {
		args.Key = key
	}

	// increment 21
	cryptoKey := os.Getenv("CryptoKey")
	if cryptoKey != "" {
		args.CryptoKey = cryptoKey
	}

	if args.Addr == "" {
		args.Addr = flagAddr
	}
	if args.LogLevel == "" {
		args.LogLevel = flagLogLevel
	}

	if args.StoreInterval == 0 {
		args.StoreInterval = flagStoreInterval
	}
	if args.StoragePath == "" {
		args.StoragePath = flagStoragePath
	}
	if !args.Restore {
		args.Restore = flagRestore
	}

	if args.DatabaseURL == "" {
		args.DatabaseURL = flagDatabaseStr
	}

	if args.Key == "" {
		args.Key = flagKey
	}

	if args.CryptoKey == "" {
		args.CryptoKey = flagCryptoKey
	}

	return args
}
