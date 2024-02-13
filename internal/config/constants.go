package config

const (
	CryptoHeader = "Crypto-Key"
	SignHeader   = "HashSHA256"
	RealIPHeader = "X-Real-IP"

	// OS enviroments
	ConfFile  = "CONFIG"
	CryptoKey = "CryptoKey"
	Address   = "ADDRESS"
	Key       = "KEY"
	LogLevel  = "LOGLEVEL"

	// Protocols
	Protocol  = "PROTO"
	ProtoGRPC = "grpc"
	ProtoHTTP = "http"

	// server only constants
	Restore       = "RESTORE"
	StoragePath   = "FILE_STORAGE_PATH"
	DBURL         = "DATABASE_DSN"
	StoreInterval = "STORE_INTERVAL"
	CIDR          = "TRUSTED_SUBNET"

	// agent only constants
	ReportInterval  = "REPORT_INTERVAL"
	RateLimit       = "RATE_LIMIT"
	PollInterval    = "POLL_INTERVAL"
	DefaultLogLevel = "INFO"
)
