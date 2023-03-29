package service

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/health"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dpkafka "github.com/ONSdigital/dp-kafka/v3"
	dphttp "github.com/ONSdigital/dp-net/v2/http"
	"github.com/ONSdigital/dp-search-reindex-tracker/config"
)

// ExternalServiceList holds the initialiser and initialisation state of external services.
type ExternalServiceList struct {
	HealthCheck    bool
	KafkaConsumers bool
	Init           Initialiser
}

// NewServiceList creates a new service list with the provided initialiser
func NewServiceList(initialiser Initialiser) *ExternalServiceList {
	return &ExternalServiceList{
		HealthCheck:    false,
		KafkaConsumers: false,
		Init:           initialiser,
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
func (e *ExternalServiceList) GetKafkaConsumers(ctx context.Context, cfg *config.Config) (reindexRequestedConsumer, reindexTaskCountsConsumer, searchDataImportedConsumer dpkafka.IConsumerGroup, err error) {
	reindexRequestedConsumer, reindexTaskCountsConsumer, searchDataImportedConsumer, err = e.Init.DoGetKafkaConsumers(ctx, &cfg.KafkaConfig)
	if err != nil {
		return
	}

	e.KafkaConsumers = true
	return
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

// GetHealthClient returns a healthclient for the provided URL
func (e *ExternalServiceList) GetHealthClient(name, url string) *health.Client {
	return e.Init.DoGetHealthClient(name, url)
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
func (e *Init) DoGetKafkaConsumers(ctx context.Context, kafkaCfg *config.KafkaConfig) (reindexRequestedConsumer, reindexTaskCountsConsumer, searchDataImportedConsumer dpkafka.IConsumerGroup, err error) {
	kafkaOffset := dpkafka.OffsetNewest
	if kafkaCfg.OffsetOldest {
		kafkaOffset = dpkafka.OffsetOldest
	}
	reindexRequestedCgConfig := &dpkafka.ConsumerGroupConfig{
		BrokerAddrs:  kafkaCfg.Brokers,
		GroupName:    kafkaCfg.ReindexRequestedGroup,
		KafkaVersion: &kafkaCfg.Version,
		NumWorkers:   &kafkaCfg.NumWorkers,
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
	reindexRequestedConsumer, err = dpkafka.NewConsumerGroup(
		ctx,
		reindexRequestedCgConfig,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	reindexTaskCountsCgConfig := &dpkafka.ConsumerGroupConfig{
		BrokerAddrs:  kafkaCfg.Brokers,
		GroupName:    kafkaCfg.ReindexTaskCountsGroup,
		KafkaVersion: &kafkaCfg.Version,
		NumWorkers:   &kafkaCfg.NumWorkers,
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
	reindexTaskCountsConsumer, err = dpkafka.NewConsumerGroup(
		ctx,
		reindexTaskCountsCgConfig,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	searchDataImportedCgConfig := &dpkafka.ConsumerGroupConfig{
		BrokerAddrs:  kafkaCfg.Brokers,
		GroupName:    kafkaCfg.SearchDataImportedGroup,
		KafkaVersion: &kafkaCfg.Version,
		NumWorkers:   &kafkaCfg.NumWorkers,
		Offset:       &kafkaOffset,
		Topic:        kafkaCfg.SearchDataImportedTopic,
	}
	if kafkaCfg.SecProtocol == config.KafkaTLSProtocolFlag {
		searchDataImportedCgConfig.SecurityConfig = dpkafka.GetSecurityConfig(
			kafkaCfg.SecCACerts,
			kafkaCfg.SecClientCert,
			kafkaCfg.SecClientKey,
			kafkaCfg.SecSkipVerify,
		)
	}
	searchDataImportedConsumer, err = dpkafka.NewConsumerGroup(
		ctx,
		searchDataImportedCgConfig,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	return reindexRequestedConsumer, reindexTaskCountsConsumer, searchDataImportedConsumer, nil
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

// DoGetHealthClient creates a new Health Client for the provided name and url
func (e *Init) DoGetHealthClient(name, url string) *health.Client {
	return health.NewClient(name, url)
}

// NewMockHTTPClient mocks HTTP Client
func NewMockHTTPClient(r *http.Response, err error) *dphttp.ClienterMock {
	return &dphttp.ClienterMock{
		SetPathsWithNoRetriesFunc: func(paths []string) {},
		GetPathsWithNoRetriesFunc: func() []string { return []string{} },
		DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
			return r, err
		},
	}
}
