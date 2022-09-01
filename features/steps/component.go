package steps

import (
	"context"
	"fmt"
	"net/http"
	"os"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-kafka/v3/kafkatest"
	"github.com/ONSdigital/dp-search-reindex-tracker/config"
	"github.com/ONSdigital/dp-search-reindex-tracker/service"
	"github.com/ONSdigital/dp-search-reindex-tracker/service/mock"
)

type SearchReindexTrackerComponent struct {
	componenttest.ErrorFeature
	serviceList                *service.ExternalServiceList
	reindexRequestedConsumer   kafka.IConsumerGroup
	reindexTaskCountsConsumer  kafka.IConsumerGroup
	searchDataImportedConsumer kafka.IConsumerGroup
	killChannel                chan os.Signal
	apiFeature                 *componenttest.APIFeature
	errorChan                  chan error
	svc                        *service.Service
	cfg                        *config.Config
	HTTPServer                 *http.Server
}

func NewSearchReindexTrackerComponent() (*SearchReindexTrackerComponent, error) {

	c := &SearchReindexTrackerComponent{
		HTTPServer: &http.Server{},
		errorChan:  make(chan error),
	}

	reindexRequestedConsumer := kafkatest.NewMessageConsumer(false)
	reindexRequestedConsumer.CheckerFunc = funcCheck
	reindexRequestedConsumer.StartFunc = func() error { return nil }
	reindexRequestedConsumer.LogErrorsFunc = func(ctx context.Context) {}
	reindexRequestedConsumer.StopFunc = func() error { return nil }
	c.reindexRequestedConsumer = reindexRequestedConsumer

	reindexTaskCountsConsumer := kafkatest.NewMessageConsumer(false)
	reindexTaskCountsConsumer.CheckerFunc = funcCheck
	reindexTaskCountsConsumer.StartFunc = func() error { return nil }
	reindexTaskCountsConsumer.LogErrorsFunc = func(ctx context.Context) {}
	reindexTaskCountsConsumer.StopFunc = func() error { return nil }
	c.reindexTaskCountsConsumer = reindexTaskCountsConsumer

	searchDataImportedConsumer := kafkatest.NewMessageConsumer(false)
	searchDataImportedConsumer.CheckerFunc = funcCheck
	searchDataImportedConsumer.StartFunc = func() error { return nil }
	searchDataImportedConsumer.LogErrorsFunc = func(ctx context.Context) {}
	searchDataImportedConsumer.StopFunc = func() error { return nil }
	c.searchDataImportedConsumer = searchDataImportedConsumer

	cfg, err := config.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get config - err: %v", err)
	}
	c.cfg = cfg

	initMock := &mock.InitialiserMock{
		DoGetKafkaConsumersFunc: c.DoGetConsumers,
		DoGetHealthCheckFunc:    c.DoGetHealthCheck,
		DoGetHTTPServerFunc:     c.DoGetHTTPServer,
	}

	c.serviceList = service.NewServiceList(initMock)

	return c, nil
}

func (c *SearchReindexTrackerComponent) Close() {
}

func (c *SearchReindexTrackerComponent) Reset() {
}

func (c *SearchReindexTrackerComponent) DoGetHealthCheck(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
	return &mock.HealthCheckerMock{
		AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
		StartFunc:    func(ctx context.Context) {},
		StopFunc:     func() {},
	}, nil
}

func (c *SearchReindexTrackerComponent) DoGetHTTPServer(bindAddr string, router http.Handler) service.HTTPServer {
	c.HTTPServer.Addr = bindAddr
	c.HTTPServer.Handler = router
	return c.HTTPServer
}

func (c *SearchReindexTrackerComponent) DoGetConsumers(ctx context.Context, kafkaCfg *config.KafkaConfig) (kafka.IConsumerGroup, kafka.IConsumerGroup, kafka.IConsumerGroup, error) {
	return c.reindexRequestedConsumer, c.reindexTaskCountsConsumer, c.searchDataImportedConsumer, nil
}

func funcCheck(ctx context.Context, state *healthcheck.CheckState) error {
	return nil
}
