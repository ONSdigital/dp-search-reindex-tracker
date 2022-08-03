package steps

import (
	"context"
	"net/http"
	"os"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-kafka/v3/kafkatest"
	dphttp "github.com/ONSdigital/dp-net/v2/http"
	"github.com/ONSdigital/dp-search-reindex-tracker/config"
	"github.com/ONSdigital/dp-search-reindex-tracker/service"
	"github.com/ONSdigital/dp-search-reindex-tracker/service/mock"
)

type Component struct {
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
}

func NewComponent() *Component {

	c := &Component{errorChan: make(chan error)}

	reindexRequestedConsumer := kafkatest.NewMessageConsumer(false)
	reindexRequestedConsumer.CheckerFunc = funcCheck
	reindexRequestedConsumer.StartFunc = func() error { return nil }
	reindexRequestedConsumer.LogErrorsFunc = func(ctx context.Context) {}
	c.reindexRequestedConsumer = reindexRequestedConsumer

	reindexTaskCountsConsumer := kafkatest.NewMessageConsumer(false)
	reindexTaskCountsConsumer.CheckerFunc = funcCheck
	reindexTaskCountsConsumer.StartFunc = func() error { return nil }
	reindexTaskCountsConsumer.LogErrorsFunc = func(ctx context.Context) {}
	c.reindexTaskCountsConsumer = reindexTaskCountsConsumer

	searchDataImportedConsumer := kafkatest.NewMessageConsumer(false)
	searchDataImportedConsumer.CheckerFunc = funcCheck
	searchDataImportedConsumer.StartFunc = func() error { return nil }
	searchDataImportedConsumer.LogErrorsFunc = func(ctx context.Context) {}
	c.searchDataImportedConsumer = searchDataImportedConsumer

	cfg, err := config.Get()
	if err != nil {
		return nil
	}

	c.cfg = cfg

	initMock := &mock.InitialiserMock{
		DoGetKafkaConsumersFunc: c.DoGetConsumers,
		DoGetHealthCheckFunc:    c.DoGetHealthCheck,
		DoGetHTTPServerFunc:     c.DoGetHTTPServer,
	}

	c.serviceList = service.NewServiceList(initMock)

	return c
}

func (c *Component) Close() {
	os.Remove(c.cfg.OutputFilePath)
}

func (c *Component) Reset() {
	os.Remove(c.cfg.OutputFilePath)
}

func (c *Component) DoGetHealthCheck(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
	return &mock.HealthCheckerMock{
		AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
		StartFunc:    func(ctx context.Context) {},
		StopFunc:     func() {},
	}, nil
}

func (c *Component) DoGetHTTPServer(bindAddr string, router http.Handler) service.HTTPServer {
	return dphttp.NewServer(bindAddr, router)
}

func (c *Component) DoGetConsumers(ctx context.Context, kafkaCfg *config.KafkaConfig) (kafka.IConsumerGroup, kafka.IConsumerGroup, kafka.IConsumerGroup, error) {
	return c.reindexRequestedConsumer, c.reindexTaskCountsConsumer, c.searchDataImportedConsumer, nil
}

func funcCheck(ctx context.Context, state *healthcheck.CheckState) error {
	return nil
}
