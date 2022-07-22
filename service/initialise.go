package service

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dpkafka "github.com/ONSdigital/dp-kafka/v3"
	dphttp "github.com/ONSdigital/dp-net/v2/http"
	"github.com/ONSdigital/dp-search-reindex-tracker/config"
)

// ExternalServiceList holds the initialiser and initialisation state of external services.
type ExternalServiceList struct {
	HealthCheck   bool
	KafkaConsumer bool
	Init          Initialiser
}

// NewServiceList creates a new service list with the provided initialiser
func NewServiceList(initialiser Initialiser) *ExternalServiceList {
	return &ExternalServiceList{
		HealthCheck:   false,
		KafkaConsumer: false,
		Init:          initialiser,
	}
}

// Init implements the Initialiser interface to initialise dependencies
type Init struct{}

// GetHTTPServer creates an http server and sets the Server flag to true
func (e *ExternalServiceList) GetHTTPServer(bindAddr string, router http.Handler) HTTPServer {
	s := e.Init.DoGetHTTPServer(bindAddr, router)
	return s
}

// GetKafkaConsumers creates three Kafka consumers and sets the consumer flag to true
func (e *ExternalServiceList) GetKafkaConsumers(ctx context.Context, cfg *config.Config) (dpkafka.IConsumerGroup, dpkafka.IConsumerGroup, error) {
	reindexRequestedConsumer, reindexTaskCountsConsumer, err := e.Init.DoGetKafkaConsumers(ctx, &cfg.KafkaConfig)
	if err != nil {
		return nil, nil, err
	}
	e.KafkaConsumer = true
	return reindexRequestedConsumer, reindexTaskCountsConsumer, nil
}

// GetHealthCheck creates a healthcheck with versionInfo and sets teh HealthCheck flag to true
func (e *ExternalServiceList) GetHealthCheck(cfg *config.Config, buildTime, gitCommit, version string) (HealthChecker, error) {
	hc, err := e.Init.DoGetHealthCheck(cfg, buildTime, gitCommit, version)
	if err != nil {
		return nil, err
	}
	e.HealthCheck = true
	return hc, nil
}

// DoGetHTTPServer creates an HTTP Server with the provided bind address and router
func (e *Init) DoGetHTTPServer(bindAddr string, router http.Handler) HTTPServer {
	s := dphttp.NewServer(bindAddr, router)
	s.HandleOSSignals = false
	return s
}

// DoGetKafkaConsumers returns three Kafka Consumer groups:
// - reindexRequestedConsumer
// - reindexTaskCountsConsumer
func (e *Init) DoGetKafkaConsumers(ctx context.Context, kafkaCfg *config.KafkaConfig) (dpkafka.IConsumerGroup, dpkafka.IConsumerGroup, error) {
	kafkaOffset := dpkafka.OffsetNewest
	if kafkaCfg.OffsetOldest {
		kafkaOffset = dpkafka.OffsetOldest
	}
	reindexRequestedCgConfig := &dpkafka.ConsumerGroupConfig{
		BrokerAddrs:  kafkaCfg.Brokers,
		GroupName:    kafkaCfg.ReindexRequestedGroup,
		KafkaVersion: &kafkaCfg.Version,
		Offset:       &kafkaOffset,
		Topic:        kafkaCfg.ReindexRequestedTopic,
	}
	if kafkaCfg.SecProtocol == config.KafkaTLSProtocolFlag {
		reindexRequestedCgConfig.SecurityConfig = dpkafka.GetSecurityConfig(
			kafkaCfg.SecCACerts,
			kafkaCfg.SecClientCert,
			kafkaCfg.SecClientKey,
			kafkaCfg.SecSkipVerify,
		)
	}
	reindexRequestedConsumer, err := dpkafka.NewConsumerGroup(
		ctx,
		reindexRequestedCgConfig,
	)
	if err != nil {
		return nil, nil, err
	}

	reindexTaskCountsCgConfig := &dpkafka.ConsumerGroupConfig{
		BrokerAddrs:  kafkaCfg.Brokers,
		GroupName:    kafkaCfg.ReindexTaskCountsGroup,
		KafkaVersion: &kafkaCfg.Version,
		Offset:       &kafkaOffset,
		Topic:        kafkaCfg.ReindexTaskCountsTopic,
	}
	if kafkaCfg.SecProtocol == config.KafkaTLSProtocolFlag {
		reindexTaskCountsCgConfig.SecurityConfig = dpkafka.GetSecurityConfig(
			kafkaCfg.SecCACerts,
			kafkaCfg.SecClientCert,
			kafkaCfg.SecClientKey,
			kafkaCfg.SecSkipVerify,
		)
	}
	reindexTaskCountsConsumer, err := dpkafka.NewConsumerGroup(
		ctx,
		reindexTaskCountsCgConfig,
	)
	if err != nil {
		return nil, nil, err
	}

	return reindexRequestedConsumer, reindexTaskCountsConsumer, nil
}

// DoGetHealthCheck creates a healthcheck with versionInfo
func (e *Init) DoGetHealthCheck(cfg *config.Config, buildTime, gitCommit, version string) (HealthChecker, error) {
	versionInfo, err := healthcheck.NewVersionInfo(buildTime, gitCommit, version)
	if err != nil {
		return nil, err
	}
	hc := healthcheck.New(versionInfo, cfg.HealthCheckCriticalTimeout, cfg.HealthCheckInterval)
	return &hc, nil
}
