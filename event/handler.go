package event

import (
	"context"

	"github.com/ONSdigital/dp-search-reindex-tracker/config"
)

//go:generate mockgen -destination mock/handler.go -package mock github.com/ONSdigital/dp-search-reindex-tracker/event Handler

// Handler represents a handler for processing a single event.
type Handler[M KafkaAvroModel] interface {
	Handle(ctx context.Context, cfg *config.Config, topicModel *M) error
}

// ReindexRequestedHandler is the handler for reindex requested messages.
type ReindexRequestedHandler struct {
}

// Handle takes a single event
func (h *ReindexRequestedHandler) Handle(ctx context.Context, cfg *config.Config, event *ReindexRequestedModel) error {
	// TO-DO: add job id to in-memory variable for global lookup e.g. map[<job_id>] bool
	// TO-DO: make a request to the Search Reindex API (via API router) to update state to in-progress for the job
	return nil
}

// ReindexTaskCountsHandler is the handler for reindex task counts messages.
type ReindexTaskCountsHandler struct {
}

// Handle takes a single event
func (h *ReindexTaskCountsHandler) Handle(ctx context.Context, cfg *config.Config, event *ReindexTaskCountsModel) error {
	return nil
}

// SearchDataImportHandler is the handler for search data import messages.
type SearchDataImportHandler struct {
}

// Handle takes a single event
func (h *SearchDataImportHandler) Handle(ctx context.Context, cfg *config.Config, event *SearchDataImportModel) error {
	return nil
}
