package event

import (
	"context"

	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-search-reindex-tracker/config"
	"github.com/ONSdigital/log.go/v2/log"
)

// ProcessMessage unmarshals the provided kafka message into an event and calls the handler.
// Handling the commit of the message is done by the dp-kafka library which will commit the message on success or failure.
func ProcessMessage[M KafkaAvroModel](ctx context.Context, cfg *config.Config, topicEvent *KafkaConsumerEvent[M], message kafka.Message) error {

	// unmarshal message
	event, err := unmarshal[M](topicEvent.Schema, message)
	if err != nil {
		logData := log.Data{
			"schema":  topicEvent.Schema,
			"message": message,
		}
		log.Error(ctx, "failed to unmarshal event", err, logData)
		return err
	}

	log.Info(ctx, "event received", log.Data{"event": event})

	// handle event
	err = topicEvent.Handler.Handle(ctx, cfg, event)
	if err != nil {
		log.Error(ctx, "failed to handle event", err)
		return err
	}

	log.Info(ctx, "message processed")

	return nil
}
