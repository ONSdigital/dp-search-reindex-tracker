package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

const KafkaTLSProtocolFlag = "TLS"

// Config represents service configuration for dp-search-reindex-tracker
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	OutputFilePath             string        `envconfig:"OUTPUT_FILE_PATH"`
	KafkaConfig                KafkaConfig
}

// KafkaConfig contains the config required to connect to Kafka
type KafkaConfig struct {
	Brokers          []string `envconfig:"KAFKA_ADDR"`
	ReindexRequestedGroup string   `envconfig:"KAFKA_REINDEX_REQUESTED_GROUP"`
	ReindexRequestedTopic string   `envconfig:"KAFKA_REINDEX_REQUESTED_TOPIC"`
	NumWorkers       int      `envconfig:"KAFKA_NUM_WORKERS"`
	OffsetOldest     bool     `envconfig:"KAFKA_OFFSET_OLDEST"`
	SecCACerts       string   `envconfig:"KAFKA_SEC_CA_CERTS"`
	SecClientCert    string   `envconfig:"KAFKA_SEC_CLIENT_CERT"`
	SecClientKey     string   `envconfig:"KAFKA_SEC_CLIENT_KEY"    json:"-"`
	SecProtocol      string   `envconfig:"KAFKA_SEC_PROTO"`
	SecSkipVerify    bool     `envconfig:"KAFKA_SEC_SKIP_VERIFY"`
	Version          string   `envconfig:"KAFKA_VERSION"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                   "localhost:28500",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		OutputFilePath:             "/tmp/helloworld.txt",
		KafkaConfig: KafkaConfig{
			Brokers:          []string{"localhost:9092"},
			ReindexRequestedGroup: "dp-search-reindex-tracker",
			ReindexRequestedTopic: "reindex-requested",
			NumWorkers:       1,
			OffsetOldest:     true,
			SecCACerts:       "",
			SecClientCert:    "",
			SecClientKey:     "",
			SecProtocol:      "",
			SecSkipVerify:    false,
			Version:          "1.0.2",
		},
	}

	return cfg, envconfig.Process("", cfg)
}
