package event

type KafkaAvroModel interface {
	HelloCalledModel | ReindexRequestedModel
}

// TODO: remove hello called example model
// HelloCalledModel provides an avro structure for a Hello Called event
type HelloCalledModel struct {
	RecipientName string `avro:"recipient_name"`
}

// ReindexRequestedModel provides an avro structure for a Reindex Requested event
type ReindexRequestedModel struct {
	JobID       string `avro:"job_id"`
	SearchIndex string `avro:"search_index"`
	TraceID     string `avro:"trace_id"`
}
