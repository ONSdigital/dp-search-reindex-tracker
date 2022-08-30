package event

import (
	"context"

	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-search-reindex-tracker/config"
	"github.com/ONSdigital/log.go/v2/log"
)

// Consume converts messages to event instances, and pass the event to the provided handler.
func Consume[M KafkaAvroModel](ctx context.Context, cfg *config.Config, topicEvent *KafkaConsumerEvent[M]) {
	// consume loop, to be executed by each worker
	var consume = func(workerID int) {
		for {
			select {
			case message, ok := <-topicEvent.ConsumerGroup.Channels().Upstream:
				if !ok {
					log.Info(ctx, "closing event consumer loop because upstream channel is closed", log.Data{"worker_id": workerID})
					return
				}
				messageCtx := context.Background()
				processMessage(messageCtx, cfg, topicEvent, message)
				message.Release()
			case <-topicEvent.ConsumerGroup.Channels().Closer:
				log.Info(ctx, "closing event consumer loop because closer channel is closed", log.Data{"worker_id": workerID})
				return
			}
		}
	}

	// workers to consume messages in parallel
	for w := 1; w <= cfg.KafkaConfig.NumWorkers; w++ {
		go consume(w)
	}
}

// processMessage unmarshals the provided kafka message into an event and calls the handler.
// After the message is handled, it is committed.
func processMessage[M KafkaAvroModel](ctx context.Context, cfg *config.Config, topicEvent *KafkaConsumerEvent[M], message kafka.Message) {

	// unmarshal - commit on failure (consuming the message again would result in the same error)
	event, err := unmarshal[M](topicEvent.Schema, message)
	if err != nil {
		logData := log.Data{
			"schema":  topicEvent.Schema,
			"message": message,
		}
		log.Error(ctx, "failed to unmarshal event", err, logData)
		message.Commit()
		return
	}

	log.Info(ctx, "event received", log.Data{"event": event})

	// handle - commit on failure (implement error handling to not commit if message needs to be consumed again)
	err = topicEvent.Handler.Handle(ctx, cfg, event)
	if err != nil {
		log.Error(ctx, "failed to handle event", err)
		message.Commit()
		return
	}

	log.Info(ctx, "event processed - committing message", log.Data{"event": event})
	message.Commit()
	log.Info(ctx, "message committed", log.Data{"event": event})
}
