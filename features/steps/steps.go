package steps

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-kafka/v3/kafkatest"
	"github.com/ONSdigital/dp-search-reindex-tracker/event"
	"github.com/ONSdigital/dp-search-reindex-tracker/service"
	"github.com/ONSdigital/go-ns/avro"
	"github.com/cucumber/godog"
	"github.com/rdumont/assistdog"
)

func (c *SearchReindexTrackerComponent) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^these reindex-requested events are consumed:$`, c.theseReindexRequestedEventsAreConsumed)
	ctx.Step(`^these reindex-task-counts events are consumed:$`, c.theseReindexTaskCountsEventsAreConsumed)
	ctx.Step(`^these search-data-import events are consumed:$`, c.theseSearchDataImportEventsAreConsumed)
	ctx.Step(`^nothing happens`, c.nothingHappens)
}

func (c *SearchReindexTrackerComponent) nothingHappens() error {
	return nil
}

func (c *SearchReindexTrackerComponent) theseReindexRequestedEventsAreConsumed(table *godog.Table) error {
	return theseEventsAreConsumed[event.ReindexRequestedModel](c, event.ReindexRequestedSchema, c.reindexRequestedConsumer, table)
}

func (c *SearchReindexTrackerComponent) theseReindexTaskCountsEventsAreConsumed(table *godog.Table) error {
	return theseEventsAreConsumed[event.ReindexTaskCountsModel](c, event.ReindexTaskCountsSchema, c.reindexTaskCountsConsumer, table)
}

func (c *SearchReindexTrackerComponent) theseSearchDataImportEventsAreConsumed(table *godog.Table) error {
	return theseEventsAreConsumed[event.SearchDataImportModel](c, event.SearchDataImportSchema, c.searchDataImportedConsumer, table)
}

func theseEventsAreConsumed[M event.KafkaAvroModel](c *SearchReindexTrackerComponent, schema *avro.Schema, consumer kafka.IConsumerGroup, table *godog.Table) error {
	ctx := context.Background()

	observationEvents, err := convertToEvents[M](table)
	if err != nil {
		return fmt.Errorf("failed to convert to events - err: %v", err)
	}

	signals := registerInterrupt()

	// run application in separate goroutine
	go func() {
		c.svc, err = service.Run(ctx, c.serviceList, "", "", "", c.errorChan)
	}()

	// consume extracted observations
	for _, e := range observationEvents {
		bytes, err := schema.Marshal(e)
		if err != nil {
			return fmt.Errorf("failed to marshal events - err: %v", err)
		}

		consumer.Channels().Upstream <- kafkatest.NewMessage(bytes, 0)
	}

	time.Sleep(300 * time.Millisecond)

	// kill application
	signals <- os.Interrupt

	return nil
}

func convertToEvents[M event.KafkaAvroModel](table *godog.Table) ([]*M, error) {
	var emptyEventModel M
	assist := assistdog.NewDefault()

	// register parser for handling int32 types as the default assist does not handle int32 types
	assist.RegisterParser(int32(0), func(raw string) (interface{}, error) {
		valInt, err := strconv.Atoi(raw)
		if err != nil {
			return nil, fmt.Errorf("failed to register int32 parser in assist - err: %v", err)
		}
		return int32(valInt), nil
	})

	events, err := assist.CreateSlice(&emptyEventModel, table)
	if err != nil {
		return nil, fmt.Errorf("failed to create events - err: %v", err)
	}
	return events.([]*M), nil
}

func registerInterrupt() chan os.Signal {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	return signals
}
