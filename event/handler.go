package event

import (
	"context"
	"fmt"
	"sync"

	searchReindexAPIModel "github.com/ONSdigital/dp-search-reindex-api/models"
	searchReindexAPIClient "github.com/ONSdigital/dp-search-reindex-api/sdk"
	"github.com/ONSdigital/dp-search-reindex-tracker/config"
	"github.com/ONSdigital/log.go/v2/log"
)

const patchReplace = "replace"

var jobsInProgress sync.Map

//go:generate mockgen -destination mock/handler.go -package mock github.com/ONSdigital/dp-search-reindex-tracker/event Handler

// Handler represents a handler for processing a single event.
type Handler[M KafkaAvroModel] interface {
	Handle(ctx context.Context, cfg *config.Config, topicModel *M) error
}

// ReindexRequestedHandler is the handler for reindex requested messages.
type ReindexRequestedHandler struct {
	SearchReindexAPIClient searchReindexAPIClient.Client
}

// Handle takes a single reindex-requested event
func (h *ReindexRequestedHandler) Handle(ctx context.Context, cfg *config.Config, event *ReindexRequestedModel) error {
	// make a patch request to the Search Reindex API to update state to in-progress for the job
	reqHeaders := searchReindexAPIClient.Headers{
		ServiceAuthToken: cfg.ServiceAuthToken,
	}
	reqPatchOp := []searchReindexAPIClient.PatchOperation{
		{
			Op:    patchReplace,
			Path:  searchReindexAPIModel.JobStatePath,
			Value: searchReindexAPIModel.JobStateInProgress,
		},
	}
	_, err := h.SearchReindexAPIClient.PatchJob(ctx, reqHeaders, event.JobID, reqPatchOp)
	if err != nil {
		log.Error(ctx, "failed to make patch request to search-reindex-api", err)
		return err
	}

	// add job id to in-memory variable for global lookup (map[<job_id>] bool)
	jobsInProgress.Store(event.JobID, true)

	return nil
}

// ReindexTaskCountsHandler is the handler for reindex task counts messages.
type ReindexTaskCountsHandler struct {
	SearchReindexAPIClient searchReindexAPIClient.Client
}

// Handle updates the number of documents associated with a task and then does a patch to the job to update the same.
func (h *ReindexTaskCountsHandler) Handle(ctx context.Context, cfg *config.Config, event *ReindexTaskCountsModel) error {
	reqHeaders := searchReindexAPIClient.Headers{
		ServiceAuthToken: cfg.ServiceAuthToken,
	}

	_, _, err := h.SearchReindexAPIClient.PostTask(ctx, reqHeaders, event.JobID, searchReindexAPIModel.TaskToCreate{
		TaskName:          event.TaskName,
		NumberOfDocuments: int(event.TaskCount),
	})
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed to create task %v associated with the job with id : %v", event.TaskName, event.JobID), err)
		return err
	}

	_, tasks, err := h.SearchReindexAPIClient.GetTasks(ctx, reqHeaders, event.JobID)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed to retrieve number of tasks associated with the job with id : %v", event.JobID), err)
		return err
	}

	totalNumberOfTasksCompleted := tasks.TotalCount + 1
	reqPatchOp := []searchReindexAPIClient.PatchOperation{
		{
			Op:    patchReplace,
			Path:  searchReindexAPIModel.JobTotalSearchDocumentsPath,
			Value: event.TaskCount,
		},
		{
			Op:    patchReplace,
			Path:  searchReindexAPIModel.JobURLExtractionCompletedStatusPath,
			Value: event.ExtractionCompleted,
		},
		{
			Op:    patchReplace,
			Path:  searchReindexAPIModel.JobNoOfTasksPath,
			Value: totalNumberOfTasksCompleted,
		},
	}
	_, err = h.SearchReindexAPIClient.PatchJob(ctx, reqHeaders, event.JobID, reqPatchOp)
	if err != nil {
		log.Error(ctx, "failed to make patch request to search-reindex-api", err)
		return err
	}

	return nil
}

// SearchDataImportHandler is the handler for search data import messages.
type SearchDataImportHandler struct {
}

// TODO: SearchDataImportHandler.Handle takes a single search-data-import event and handles it
func (h *SearchDataImportHandler) Handle(ctx context.Context, cfg *config.Config, event *SearchDataImportModel) error {
	return nil
}
