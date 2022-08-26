package event

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/ONSdigital/dp-search-reindex-tracker/config"
	"github.com/ONSdigital/log.go/v2/log"
)

//go:generate mockgen -destination mock/handler.go -package mock github.com/ONSdigital/dp-search-reindex-tracker/event Handler

// TODO: remove or replace hello called logic with app specific
// Handler represents a handler for processing a single event.
type Handler[kafkaTopicModel KafkaAvroModel] interface {
	Handle(ctx context.Context, cfg *config.Config, topicModel *kafkaTopicModel) error
}

// TODO: remove hello called example handler
// HelloCalledHandler ...
type HelloCalledHandler struct {
}

// Handle takes a single event.
func (h *HelloCalledHandler) Handle(ctx context.Context, cfg *config.Config, event *HelloCalledModel) (err error) {
	logData := log.Data{
		"event": event,
	}
	log.Info(ctx, "event handler called", logData)

	greeting := fmt.Sprintf("Hello, %s!", event.RecipientName)
	err = ioutil.WriteFile(cfg.OutputFilePath, []byte(greeting), 0644)
	if err != nil {
		return err
	}

	logData["greeting"] = greeting
	log.Info(ctx, "hello world example handler called successfully", logData)
	log.Info(ctx, "event successfully handled", logData)

	return nil
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
