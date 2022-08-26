package event

import (
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/go-ns/avro"
)

// unmarshal converts a event instance to []byte.
func unmarshal[topicModel KafkaAvroModel](topicSchema *avro.Schema, message kafka.Message) (*topicModel, error) {
	var event topicModel
	err := topicSchema.Unmarshal(message.GetData(), &event)
	return &event, err
}

// TODO: remove or replace hello called structure and model with app specific
var helloCalledSchema = `{
  "type": "record",
  "name": "hello-called",
  "fields": [
    {"name": "recipient_name", "type": "string", "default": ""}
  ]
}`

// HelloCalledSchema is the Avro schema for Hello Called messages.
var HelloCalledSchema = &avro.Schema{
	Definition: helloCalledSchema,
}

var reindexRequestedSchema = `{
	"type": "record",
	"name": "reindex-requested",
	"fields": [
		{"name": "job_id", "type": "string", "default": ""},
		{"name": "search_index", "type": "string", "default": ""},
		{"name": "trace_id", "type": "string", "default": ""}
	]
}`

// ReindexRequestedSchema is the Avro schema for Reindex Requested messages.
var ReindexRequestedSchema = &avro.Schema{
	Definition: reindexRequestedSchema,
}
