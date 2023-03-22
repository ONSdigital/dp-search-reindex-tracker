package event

import (
	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/go-ns/avro"
)

var reindexRequested = `{
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
	Definition: reindexRequested,
}

var reindexTaskCounts = `{
	"type": "record",
	"name": "reindex-task-counts",
	"fields": "fields": [
    {"name": "job_id", "type": "string", "default": ""},
    {"name": "task", "type": "string", "default": ""},
    {"name": "extraction_completed", "type": "boolean", "default": false},
    {"name": "count", "type": "string", "default":"0"}
  ]
}`

// ReindexTaskCountsSchema is the Avro schema for Reindex Task Counts messages.
var ReindexTaskCountsSchema = &avro.Schema{
	Definition: reindexTaskCounts,
}

var searchDataImport = `{
	"type": "record",
	"name": "search-data-import",
	"fields": [
	  {"name": "uid", "type": "string", "default": ""},
	  {"name": "uri", "type": "string", "default": ""},
	  {"name": "data_type", "type": "string", "default": ""},
	  {"name": "job_id", "type": "string", "default": ""},
	  {"name": "search_index", "type": "string", "default": ""},
	  {"name": "cdid", "type": "string", "default": ""},
	  {"name": "dataset_id", "type": "string", "default": ""},
	  {"name": "edition", "type": "string", "default": ""},
	  {"name": "keywords", "type": {"type":"array","items":"string"}},
	  {"name": "meta_description", "type": "string", "default": ""},
	  {"name": "release_date", "type": "string", "default": ""},
	  {"name": "summary", "type": "string", "default": ""},
	  {"name": "title", "type": "string", "default": ""},
	  {"name": "topics", "type": {"type":"array","items":"string"}},
	  {"name": "trace_id", "type": "string", "default": ""},
	  {"name": "cancelled", "type": "boolean", "default": false},
	  {"name": "finalised", "type": "boolean", "default": false},
	  {"name": "published", "type": "boolean", "default": false},
	  {"name": "language", "type": "string", "default": ""},
	  {"name": "survey",   "type": "string", "default": ""},
	  {"name": "canonical_topic",   "type": "string", "default": ""},
	  {"name": "date_changes", "type": {"type":"array","items":{
	   "name": "ReleaseDateDetails",
	   "type" : "record",
	   "fields" : [
		{"name": "change_notice", "type": "string", "default": ""},
		{"name": "previous_date", "type": "string", "default": ""}
	  ]}}},
	  {"name": "provisional_date", "type": "string", "default": ""}
	]
  }`

// SearchDataImportSchema the Avro schema for search-data-import messages.
var SearchDataImportSchema = &avro.Schema{
	Definition: searchDataImport,
}

// unmarshal converts a event instance to []byte.
func unmarshal[M KafkaAvroModel](topicSchema *avro.Schema, message kafka.Message) (*M, error) {
	var event M
	err := topicSchema.Unmarshal(message.GetData(), &event)
	return &event, err
}
