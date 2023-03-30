package event

type KafkaAvroModel interface {
	ReindexRequestedModel | SearchDataImportModel | ReindexTaskCountsModel
}

// ReindexRequestedModel provides an avro structure for a Reindex Requested event
type ReindexRequestedModel struct {
	JobID       string `avro:"job_id"`
	SearchIndex string `avro:"search_index"`
	TraceID     string `avro:"trace_id"`
}

// ReindexTaskCountsModel provides an avro structure for a Reindex Task Counts event
type ReindexTaskCountsModel struct {
	JobID               string `avro:"job_id"`
	TaskName            string `avro:"task"`
	ExtractionCompleted bool   `avro:"extraction_completed"`
	TaskCount           int32  `avro:"count"`
}

// SearchDataImport provides event data for a search data import
type SearchDataImportModel struct {
	UID             string               `avro:"uid"`
	URI             string               `avro:"uri"`
	Edition         string               `avro:"edition"`
	DataType        string               `avro:"data_type"`
	JobID           string               `avro:"job_id"`
	SearchIndex     string               `avro:"search_index"`
	CDID            string               `avro:"cdid"`
	DatasetID       string               `avro:"dataset_id"`
	Keywords        []string             `avro:"keywords"`
	MetaDescription string               `avro:"meta_description"`
	ReleaseDate     string               `avro:"release_date"`
	Summary         string               `avro:"summary"`
	Title           string               `avro:"title"`
	Topics          []string             `avro:"topics"`
	TraceID         string               `avro:"trace_id"`
	DateChanges     []ReleaseDateDetails `avro:"date_changes"`
	Cancelled       bool                 `avro:"cancelled"`
	Finalised       bool                 `avro:"finalised"`
	ProvisionalDate string               `avro:"provisional_date"`
	CanonicalTopic  string               `avro:"canonical_topic"`
	Published       bool                 `avro:"published"`
	Language        string               `avro:"language"`
	Survey          string               `avro:"survey"`
}

// ReleaseDateChange represent a date change of a release
type ReleaseDateDetails struct {
	ChangeNotice string `avro:"change_notice"`
	Date         string `avro:"previous_date"`
}
