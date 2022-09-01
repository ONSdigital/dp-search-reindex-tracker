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

// GetReindexRequested returns information needed to handle the reindex-requested kafka messages
// which includes its corresponding handler, schema and other options such as the consumer group
func GetReindexRequested(options *KafkaEventOptions) *KafkaConsumerEvent[ReindexRequestedModel] {
	topic := KafkaConsumerEvent[ReindexRequestedModel]{
		Handler: &ReindexRequestedHandler{},
		Schema:  ReindexRequestedSchema,
	}

	topic.ConsumerGroup = options.ConsumerGroup
	return &topic
}

// GetReindexTaskCounts returns information needed to handle the reindex-task-counts kafka messages
// which includes its corresponding handler, schema and other options such as the consumer group
func GetReindexTaskCounts(options *KafkaEventOptions) *KafkaConsumerEvent[ReindexTaskCountsModel] {
	topic := KafkaConsumerEvent[ReindexTaskCountsModel]{
		Handler: &ReindexTaskCountsHandler{},
		Schema:  ReindexTaskCountsSchema,
	}

	topic.ConsumerGroup = options.ConsumerGroup
	return &topic
}

// GetSearchDataImport returns information needed to handle the search-data-import kafka messages
// which includes its corresponding handler, schema and other options such as the consumer group
func GetSearchDataImport(options *KafkaEventOptions) *KafkaConsumerEvent[SearchDataImportModel] {
	topic := KafkaConsumerEvent[SearchDataImportModel]{
		Handler: &SearchDataImportHandler{},
		Schema:  SearchDataImportSchema,
	}

	topic.ConsumerGroup = options.ConsumerGroup
	return &topic
}
