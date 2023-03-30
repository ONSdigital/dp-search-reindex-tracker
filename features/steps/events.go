package steps

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-kafka/v3/kafkatest"
	"github.com/ONSdigital/dp-search-reindex-tracker/event"
	"github.com/ONSdigital/go-ns/avro"
	"github.com/cucumber/godog"
	"github.com/rdumont/assistdog"
)

func theseEventsAreConsumed[M event.KafkaAvroModel](c *SearchReindexTrackerComponent, schema *avro.Schema, consumer kafka.IConsumerGroup, table *godog.Table) error {
	observationEvents, err := convertToEvents[M](table)
	if err != nil {
		return fmt.Errorf("failed to convert to events - err: %v", err)
	}

	// consume extracted observations
	for _, e := range observationEvents {
		bytes, err := schema.Marshal(e)
		if err != nil {
			return fmt.Errorf("failed to marshal events - err: %v", err)
		}

		consumer.Channels().Upstream <- kafkatest.NewMessage(bytes, 0)
	}

	time.Sleep(300 * time.Millisecond)

	return nil
}

func convertToEvents[M event.KafkaAvroModel](table *godog.Table) ([]*M, error) {
	var emptyEventModel M
	assist := assistdog.NewDefault()

	// register parser for handling int32 types as the default assist does not handle int32 types
	assist.RegisterParser(int32(0), func(raw string) (interface{}, error) {
		valInt, err := strconv.ParseInt(raw, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("failed to register int32 parser in assist - err: %v", err)
		}

		return int32(valInt), nil
	})

	// register parser for handling boolean types as the default assist does not handle boolean types
	assist.RegisterParser(false, func(raw string) (interface{}, error) {
		valBool, err := strconv.ParseBool(raw)
		if err != nil {
			return nil, fmt.Errorf("failed to register boolean parser in assist - err: %v", err)
		}
		return valBool, nil
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
