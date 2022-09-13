package steps

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ONSdigital/dp-api-clients-go/v2/health"
	componentTest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-kafka/v3/kafkatest"
	"github.com/ONSdigital/dp-search-reindex-tracker/config"
	"github.com/ONSdigital/dp-search-reindex-tracker/service"
	"github.com/ONSdigital/dp-search-reindex-tracker/service/mock"
)

const (
	gitCommitHash = "3t7e5s1t4272646ef477f8ed755"
	appVersion    = "v1.2.3"
)

type SearchReindexTrackerComponent struct {
	componentTest.ErrorFeature
	serviceList                *service.ExternalServiceList
	fakeAPIRouter              *FakeAPI
	reindexRequestedConsumer   kafka.IConsumerGroup
	reindexTaskCountsConsumer  kafka.IConsumerGroup
	searchDataImportedConsumer kafka.IConsumerGroup
	apiFeature                 *componentTest.APIFeature
	errorChan                  chan error
	svc                        *service.Service
	ctx                        context.Context
	cfg                        *config.Config
	HTTPServer                 *http.Server
	startTime                  time.Time
}

func NewSearchReindexTrackerComponent() (*SearchReindexTrackerComponent, error) {
	c := &SearchReindexTrackerComponent{
		HTTPServer: &http.Server{},
		errorChan:  make(chan error),
	}

	c.ctx = context.Background()

	cfg, err := config.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get config - err: %v", err)
	}
	c.cfg = cfg

	c.cfg.HealthCheckInterval = 1 * time.Second
	c.cfg.HealthCheckCriticalTimeout = 3 * time.Second

	c.fakeAPIRouter = NewFakeAPI()
	c.cfg.APIRouterURL = c.fakeAPIRouter.fakeHTTP.ResolveURL("")
	c.fakeAPIRouter.healthRequest = c.fakeAPIRouter.fakeHTTP.NewHandler().Get("/health")
	c.fakeAPIRouter.healthRequest.CustomHandle = statusHandle(200)

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

	initMock := &mock.InitialiserMock{
		DoGetKafkaConsumersFunc: c.DoGetConsumers,
		DoGetHealthCheckFunc:    c.DoGetHealthCheck,
		DoGetHealthClientFunc:   c.DoGetHealthClient,
		DoGetHTTPServerFunc:     c.DoGetHTTPServer,
	}

	c.serviceList = service.NewServiceList(initMock)

	// run application in separate goroutine
	c.startTime = time.Now()
	c.svc, err = service.Run(c.ctx, c.serviceList, "", "", "", c.errorChan)
	if err != nil {
		return nil, fmt.Errorf("failed to run component service - err: %v", err)
	}

	return c, nil
}

func (c *SearchReindexTrackerComponent) Close() {
	// kill application
	signals := registerInterrupt()
	signals <- os.Interrupt
}

func (c *SearchReindexTrackerComponent) Reset() {
}

// InitAPIFeature initialises the ApiFeature that's contained within a specific JobsFeature.
func (c *SearchReindexTrackerComponent) InitAPIFeature() *componentTest.APIFeature {
	c.apiFeature = componentTest.NewAPIFeature(c.InitialiseService)

	return c.apiFeature
}

// InitialiseService returns the http.Handler that's contained within the component.
func (c *SearchReindexTrackerComponent) InitialiseService() (http.Handler, error) {
	return c.HTTPServer.Handler, nil
}

func (c *SearchReindexTrackerComponent) DoGetHealthCheck(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
	componentBuildTime := strconv.Itoa(int(time.Now().Unix()))
	versionInfo, err := healthcheck.NewVersionInfo(componentBuildTime, gitCommitHash, appVersion)
	if err != nil {
		return nil, err
	}
	hc := healthcheck.New(versionInfo, cfg.HealthCheckCriticalTimeout, cfg.HealthCheckInterval)
	return &hc, nil
}

func (c *SearchReindexTrackerComponent) DoGetHealthClient(name, url string) *health.Client {
	if name == "" || url == "" {
		return nil
	}

	return &health.Client{
		URL:    url,
		Name:   name,
		Client: c.fakeAPIRouter.getMockAPIHTTPClient(),
	}
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
	return state.Update("OK", "kafka consumer group is healthy", 200)
}
