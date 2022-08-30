package event

import (
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/go-ns/avro"
)

type KafkaConsumerEvent[M KafkaAvroModel] struct {
	ConsumerGroup kafka.IConsumerGroup
	Handler       Handler[M]
	Schema        *avro.Schema
}

type KafkaEventOptions struct {
	ConsumerGroup kafka.IConsumerGroup
}

func GetHelloCalled(options *KafkaEventOptions) *KafkaConsumerEvent[HelloCalledModel] {
	topic := KafkaConsumerEvent[HelloCalledModel]{
		Handler: &HelloCalledHandler{},
		Schema:  HelloCalledSchema,
	}

	topic.ConsumerGroup = options.ConsumerGroup
	return &topic
}

func GetReindexRequested(options *KafkaEventOptions) *KafkaConsumerEvent[ReindexRequestedModel] {
	topic := KafkaConsumerEvent[ReindexRequestedModel]{
		Handler: &ReindexRequestedHandler{},
		Schema:  ReindexRequestedSchema,
	}

	topic.ConsumerGroup = options.ConsumerGroup
	return &topic
}
