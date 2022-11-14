package event

import (
	"fmt"

	kafka "github.com/ONSdigital/dp-kafka/v3"
	searchReindexAPIClient "github.com/ONSdigital/dp-search-reindex-api/sdk"
	"github.com/ONSdigital/go-ns/avro"
)

type KafkaConsumerEvent[M KafkaAvroModel] struct {
	ConsumerGroup kafka.IConsumerGroup
	Handler       Handler[M]
	Schema        *avro.Schema
}

type KafkaEventOptions struct {
	ConsumerGroup          kafka.IConsumerGroup
	SearchReindexAPIClient searchReindexAPIClient.Client
}

func (options *KafkaEventOptions) validate() error {
	if options == nil {
		return fmt.Errorf("options for event not provided")
	}

	if options.ConsumerGroup == nil {
		return fmt.Errorf("consumer group not provided")
	}

	return nil
}

// GetReindexRequested returns information needed to handle the reindex-requested kafka messages
// which includes its corresponding handler, schema and other options such as the consumer group
func GetReindexRequested(options *KafkaEventOptions) (*KafkaConsumerEvent[ReindexRequestedModel], error) {
	err := options.validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate options given for reindex requested event")
	}

	if options.SearchReindexAPIClient == nil {
		return nil, fmt.Errorf("search reindex api client not provided")
	}

	topic := KafkaConsumerEvent[ReindexRequestedModel]{
		Handler: &ReindexRequestedHandler{
			SearchReindexAPIClient: options.SearchReindexAPIClient,
		},
		Schema: ReindexRequestedSchema,
	}

	topic.ConsumerGroup = options.ConsumerGroup
	return &topic, nil
}

// GetReindexTaskCounts returns information needed to handle the reindex-task-counts kafka messages
// which includes its corresponding handler, schema and other options such as the consumer group
func GetReindexTaskCounts(options *KafkaEventOptions) (*KafkaConsumerEvent[ReindexTaskCountsModel], error) {
	err := options.validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate options given for reindex task counts event")
	}

	topic := KafkaConsumerEvent[ReindexTaskCountsModel]{
		Handler: &ReindexTaskCountsHandler{
			SearchReindexAPIClient: options.SearchReindexAPIClient,
		},
		Schema: ReindexTaskCountsSchema,
	}

	topic.ConsumerGroup = options.ConsumerGroup
	return &topic, nil
}

// GetSearchDataImport returns information needed to handle the search-data-import kafka messages
// which includes its corresponding handler, schema and other options such as the consumer group
func GetSearchDataImport(options *KafkaEventOptions) (*KafkaConsumerEvent[SearchDataImportModel], error) {
	err := options.validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate options given for search data import event")
	}

	topic := KafkaConsumerEvent[SearchDataImportModel]{
		Handler: &SearchDataImportHandler{},
		Schema:  SearchDataImportSchema,
	}

	topic.ConsumerGroup = options.ConsumerGroup
	return &topic, nil
}
