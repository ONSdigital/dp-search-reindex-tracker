package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

const KafkaTLSProtocolFlag = "TLS"

// Config represents service configuration for dp-search-reindex-tracker
type Config struct {
	APIRouterURL               string        `envconfig:"API_ROUTER_URL"`
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	ServiceAuthToken           string        `envconfig:"SERVICE_AUTH_TOKEN"   json:"-"`
	KafkaConfig                KafkaConfig
}

// KafkaConfig contains the config required to connect to Kafka
type KafkaConfig struct {
	Brokers                 []string `envconfig:"KAFKA_ADDR"`
	NumWorkers              int      `envconfig:"KAFKA_NUM_WORKERS"`
	OffsetOldest            bool     `envconfig:"KAFKA_OFFSET_OLDEST"`
	ReindexRequestedGroup   string   `envconfig:"KAFKA_REINDEX_REQUESTED_GROUP"`
	ReindexRequestedTopic   string   `envconfig:"KAFKA_REINDEX_REQUESTED_TOPIC"`
	ReindexTaskCountsGroup  string   `envconfig:"KAFKA_REINDEX_TASK_COUNTS_GROUP"`
	ReindexTaskCountsTopic  string   `envconfig:"KAFKA_REINDEX_TASK_COUNTS_TOPIC"`
	SearchDataImportedGroup string   `envconfig:"KAFKA_SEARCH_DATA_IMPORTED_GROUP"`
	SearchDataImportedTopic string   `envconfig:"KAFKA_SEARCH_DATA_IMPORTED_TOPIC"`
	SecCACerts              string   `envconfig:"KAFKA_SEC_CA_CERTS"`
	SecClientCert           string   `envconfig:"KAFKA_SEC_CLIENT_CERT"`
	SecClientKey            string   `envconfig:"KAFKA_SEC_CLIENT_KEY"    json:"-"`
	SecProtocol             string   `envconfig:"KAFKA_SEC_PROTO"`
	SecSkipVerify           bool     `envconfig:"KAFKA_SEC_SKIP_VERIFY"`
	Version                 string   `envconfig:"KAFKA_VERSION"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		APIRouterURL:               "http://localhost:23200/v1",
		BindAddr:                   "localhost:28500",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		ServiceAuthToken:           "",
		KafkaConfig: KafkaConfig{
			Brokers:                 []string{"localhost:9092"},
			NumWorkers:              1,
			OffsetOldest:            true,
			ReindexRequestedGroup:   "dp-search-reindex-tracker",
			ReindexRequestedTopic:   "reindex-requested",
			ReindexTaskCountsGroup:  "dp-search-reindex-tracker",
			ReindexTaskCountsTopic:  "reindex-task-counts",
			SearchDataImportedGroup: "dp-search-reindex-tracker",
			SearchDataImportedTopic: "search-data-imported",
			SecCACerts:              "",
			SecClientCert:           "",
			SecClientKey:            "",
			SecProtocol:             "",
			SecSkipVerify:           false,
			Version:                 "1.0.2",
		},
	}

	return cfg, envconfig.Process("", cfg)
}
