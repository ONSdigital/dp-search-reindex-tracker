package service

import (
	"context"
	"time"

	"github.com/ONSdigital/dp-api-clients-go/v2/health"
	kafka "github.com/ONSdigital/dp-kafka/v3"
	searchReindexAPI "github.com/ONSdigital/dp-search-reindex-api/sdk/v1"
	"github.com/ONSdigital/dp-search-reindex-tracker/config"
	"github.com/ONSdigital/dp-search-reindex-tracker/event"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// Service contains all the configs, server and clients to run the event handler service
type Service struct {
	server                     HTTPServer
	router                     *mux.Router
	serviceList                *ExternalServiceList
	healthCheck                HealthChecker
	reindexRequestedConsumer   kafka.IConsumerGroup
	reindexTaskCountsConsumer  kafka.IConsumerGroup
	searchDataImportedConsumer kafka.IConsumerGroup
	shutdownTimeout            time.Duration
}

// Run the service
func Run(ctx context.Context, serviceList *ExternalServiceList, buildTime, gitCommit, version string, svcErrors chan error) (*Service, error) {
	log.Info(ctx, "running service")

	// Read config
	cfg, err := config.Get()
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve service configuration")
	}
	log.Info(ctx, "got service configuration", log.Data{"config": cfg})

	// Get HTTP Server with collectionID checkHeader middleware
	r := mux.NewRouter()
	s := serviceList.GetHTTPServer(cfg.BindAddr, r)

	// Get health client for api router
	routerHealthClient := serviceList.GetHealthClient("api-router", cfg.APIRouterURL)

	// Get search reindex API client
	searchReindexAPIClient, err := searchReindexAPI.NewWithHealthClient(cfg.ServiceAuthToken, routerHealthClient)
	if err != nil {
		log.Fatal(ctx, "failed to get search reindex api client", err)
		return nil, err
	}

	// Get Kafka consumers
	reindexRequestedConsumer, reindexTaskCountsConsumer, searchDataImportedConsumer, err := serviceList.GetKafkaConsumers(ctx, cfg)
	if err != nil {
		log.Fatal(ctx, "failed to initialise kafka consumers", err)
		return nil, err
	}

	// Get 'Reindex Requested' kafka information for handling
	reindexRequestedEventOptions := &event.KafkaEventOptions{
		ConsumerGroup:          reindexRequestedConsumer,
		SearchReindexAPIClient: searchReindexAPIClient,
	}
	reindexRequestedEvent, err := event.GetReindexRequested(reindexRequestedEventOptions)
	if err != nil {
		log.Fatal(ctx, "failed to get reindex requested event", err)
		return nil, err
	}

	// Get 'Reindex Task Counts' kafka information for handling
	reindexTaskCountsEventOptions := &event.KafkaEventOptions{
		ConsumerGroup: reindexTaskCountsConsumer,
	}
	reindexTaskCountsEvent, err := event.GetReindexTaskCounts(reindexTaskCountsEventOptions)
	if err != nil {
		log.Fatal(ctx, "failed to get reindex task counts event", err)
		return nil, err
	}

	// Get 'Search Data Imported' kafka information for handling
	searchDataImportedEventOptions := &event.KafkaEventOptions{
		ConsumerGroup: searchDataImportedConsumer,
	}
	searchDataImportedEvent, err := event.GetSearchDataImport(searchDataImportedEventOptions)
	if err != nil {
		log.Fatal(ctx, "failed to get search data imported event", err)
		return nil, err
	}

	// Kafka consumer logging go routine
	reindexRequestedConsumer.LogErrors(ctx)
	reindexTaskCountsConsumer.LogErrors(ctx)
	searchDataImportedConsumer.LogErrors(ctx)

	// Kafka consuming go routine
	if err = reindexRequestedConsumer.RegisterHandler(ctx, func(ctx context.Context, workerID int, msg kafka.Message) error {
		return event.ProcessMessage(ctx, cfg, reindexRequestedEvent, msg)
	}); err != nil {
		log.Fatal(ctx, "error registering the reindex requested consumer", err)
		return nil, err
	}

	if err = reindexTaskCountsConsumer.RegisterHandler(ctx, func(ctx context.Context, workerID int, msg kafka.Message) error {
		return event.ProcessMessage(ctx, cfg, reindexTaskCountsEvent, msg)
	}); err != nil {
		log.Fatal(ctx, "error registering the reindex task counts consumer", err)
		return nil, err
	}

	if err = searchDataImportedConsumer.RegisterHandler(ctx, func(ctx context.Context, workerID int, msg kafka.Message) error {
		return event.ProcessMessage(ctx, cfg, searchDataImportedEvent, msg)
	}); err != nil {
		log.Fatal(ctx, "error registering the search data imported consumer", err)
		return nil, err
	}

	// Start kafka consumers
	if consumerStartErr := reindexRequestedConsumer.Start(); consumerStartErr != nil {
		log.Fatal(ctx, "error starting the reindex requested consumer", consumerStartErr)
		return nil, consumerStartErr
	}
	if consumerStartErr := reindexTaskCountsConsumer.Start(); consumerStartErr != nil {
		log.Fatal(ctx, "error starting the reindex task counts consumer", consumerStartErr)
		return nil, consumerStartErr
	}
	if consumerStartErr := searchDataImportedConsumer.Start(); consumerStartErr != nil {
		log.Fatal(ctx, "error starting the search data imported consumer", consumerStartErr)
		return nil, consumerStartErr
	}

	// Get HealthCheck
	hc, err := serviceList.GetHealthCheck(cfg, buildTime, gitCommit, version)
	if err != nil {
		log.Fatal(ctx, "could not instantiate healthcheck", err)
		return nil, err
	}

	if err := registerCheckers(ctx, hc, routerHealthClient, reindexRequestedConsumer, reindexTaskCountsConsumer, searchDataImportedConsumer); err != nil {
		return nil, errors.Wrap(err, "unable to register checkers")
	}

	r.StrictSlash(true).Path("/health").HandlerFunc(hc.Handler)
	hc.Start(ctx)

	// Run the http server in a new go-routine
	go func() {
		if err := s.ListenAndServe(); err != nil {
			svcErrors <- errors.Wrap(err, "failure in http listen and serve")
		}
	}()

	return &Service{
		server:                     s,
		router:                     r,
		serviceList:                serviceList,
		healthCheck:                hc,
		reindexRequestedConsumer:   reindexRequestedConsumer,
		reindexTaskCountsConsumer:  reindexTaskCountsConsumer,
		searchDataImportedConsumer: searchDataImportedConsumer,
		shutdownTimeout:            cfg.GracefulShutdownTimeout,
	}, nil
}

// Close gracefully shuts the service down in the required order, with timeout
func (svc *Service) Close(ctx context.Context) error {
	timeout := svc.shutdownTimeout
	log.Info(ctx, "commencing graceful shutdown", log.Data{"graceful_shutdown_timeout": timeout})
	ctx, cancel := context.WithTimeout(ctx, timeout)

	// track shutdown gracefully closes up
	var gracefulShutdown bool

	go func() {
		defer cancel()
		var hasShutdownError bool

		// stop healthcheck, as it depends on everything else
		if svc.serviceList.HealthCheck {
			svc.healthCheck.Stop()
		}

		// If reindexRequestedConsumer exists, stop listening to it.
		// This will automatically stop the event consumer loops and no more messages will be processed.
		// The kafka reindexRequestedConsumer will be closed after the service shuts down.
		if svc.serviceList.KafkaConsumers {
			log.Info(ctx, "stopping reindex requested consumer listener")
			if err := svc.reindexRequestedConsumer.Stop(); err != nil {
				log.Error(ctx, "error stopping reindex requested consumer listener", err)
				hasShutdownError = true
			}
			log.Info(ctx, "stopped reindex requested consumer listener")

			log.Info(ctx, "stopping reindex task counts consumer listener")
			if err := svc.reindexTaskCountsConsumer.Stop(); err != nil {
				log.Error(ctx, "error stopping reindex task counts consumer listener", err)
				hasShutdownError = true
			}
			log.Info(ctx, "stopped reindex task counts consumer listener")

			log.Info(ctx, "stopping search data imported consumer listener")
			if err := svc.searchDataImportedConsumer.Stop(); err != nil {
				log.Error(ctx, "error stopping search data imported consumer listener", err)
				hasShutdownError = true
			}
			log.Info(ctx, "stopped search data imported consumer listener")
		}

		// stop any incoming requests before closing any outbound connections
		if err := svc.server.Shutdown(ctx); err != nil {
			log.Error(ctx, "failed to shutdown http server", err)
			hasShutdownError = true
		}

		// If kafka consumer exists, close it.
		if svc.serviceList.KafkaConsumers {
			log.Info(ctx, "closing reindex requested consumer")
			if err := svc.reindexRequestedConsumer.Close(ctx); err != nil {
				log.Error(ctx, "error closing reindex requested consumer", err)
				hasShutdownError = true
			}
			log.Info(ctx, "closed reindex requested consumer")

			log.Info(ctx, "closing reindex task counts consumer")
			if err := svc.reindexTaskCountsConsumer.Close(ctx); err != nil {
				log.Error(ctx, "error closing reindex task counts consumer", err)
				hasShutdownError = true
			}
			log.Info(ctx, "closed reindex task counts consumer")

			log.Info(ctx, "closing search data imported consumer")
			if err := svc.searchDataImportedConsumer.Close(ctx); err != nil {
				log.Error(ctx, "error closing search data imported consumer", err)
				hasShutdownError = true
			}
			log.Info(ctx, "closed search data imported consumer")
		}

		if !hasShutdownError {
			gracefulShutdown = true
		}
	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-ctx.Done()

	if !gracefulShutdown {
		err := errors.New("failed to shutdown gracefully")
		log.Error(ctx, "failed to shutdown gracefully ", err)
		return err
	}

	log.Info(ctx, "graceful shutdown was successful")
	return nil
}

func registerCheckers(ctx context.Context,
	hc HealthChecker, routerHealthClient *health.Client,
	reindexRequestedConsumer, reindexTaskCountsConsumer, searchDataImportedConsumer kafka.IConsumerGroup) (err error) {
	hasErrors := false

	if err = hc.AddCheck("API router", routerHealthClient.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "failed to add API router health checker", err)
	}

	if err := hc.AddCheck("Reindex requested consumer", reindexRequestedConsumer.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "error adding check for reindex requested consumer", err)
	}

	if err := hc.AddCheck("Reindex task counts consumer", reindexTaskCountsConsumer.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "error adding check for reindex task counts consumer", err)
	}

	if err := hc.AddCheck("Search data imported consumer", searchDataImportedConsumer.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "error adding check for search data imported consumer", err)
	}

	if hasErrors {
		return errors.New("Error(s) registering checkers for healthcheck")
	}
	return nil
}
