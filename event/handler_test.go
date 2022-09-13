package event

import (
	"context"
	"errors"
	"testing"

	"github.com/ONSdigital/dp-search-reindex-api/sdk"
	searchReindexAPIMock "github.com/ONSdigital/dp-search-reindex-api/sdk/mocks"
	"github.com/ONSdigital/dp-search-reindex-tracker/config"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	testReindexRequestEventOK = &ReindexRequestedModel{
		JobID:       "1234",
		SearchIndex: "testing",
		TraceID:     "5678",
	}
	errUnexpected = errors.New("unexpected error")
)

func TestReindexRequestedHandler_Handle(t *testing.T) {
	ctx := context.Background()

	cfg, err := config.Get()
	if err != nil {
		t.Errorf("failed to get config - err: %v", err)
	}

	Convey("Given patch request to search reindex api is successful", t, func() {
		searchReindexAPIMock := &searchReindexAPIMock.ClientMock{
			PatchJobFunc: func(ctx context.Context, headers sdk.Headers, jobID string, body []sdk.PatchOperation) (*sdk.RespHeaders, error) {
				return nil, nil
			},
		}

		reindexRequestedHandler := &ReindexRequestedHandler{
			SearchReindexAPIClient: searchReindexAPIMock,
		}

		Convey("When ReindexRequestedHandler.Handle is called", func() {
			err := reindexRequestedHandler.Handle(ctx, cfg, testReindexRequestEventOK)

			Convey("Then job id should be added to jobsInProgress for global lookup", func() {
				val, ok := jobsInProgress.Load(testReindexRequestEventOK.JobID)
				So(val, ShouldBeTrue)
				So(ok, ShouldBeTrue)

				Convey("And no error should be returned for patch request", func() {
					So(err, ShouldBeNil)

					// Reset jobsInProgress map
					jobsInProgress.Delete(testReindexRequestEventOK.JobID)
				})
			})
		})
	})

	Convey("Given patch request to search reindex api isn't successful", t, func() {
		searchReindexAPIMock := &searchReindexAPIMock.ClientMock{
			PatchJobFunc: func(ctx context.Context, headers sdk.Headers, jobID string, body []sdk.PatchOperation) (*sdk.RespHeaders, error) {
				return nil, errUnexpected
			},
		}

		reindexRequestedHandler := &ReindexRequestedHandler{
			SearchReindexAPIClient: searchReindexAPIMock,
		}

		Convey("When ReindexRequestedHandler.Handle is called", func() {
			err := reindexRequestedHandler.Handle(ctx, cfg, testReindexRequestEventOK)

			Convey("Then job id should not be added to jobsInProgress for global lookup", func() {
				_, ok := jobsInProgress.Load(testReindexRequestEventOK.JobID)
				So(ok, ShouldBeFalse)

				Convey("And an error should be returned for patch request", func() {
					So(err, ShouldNotBeNil)
				})
			})
		})
	})
}

// TODO - complete unit test once ReindexTaskCountsHandler.Handle has been implemented
func TestReindexTaskCountsHandler_Handle(t *testing.T) {
	Convey("Given a successful event handler, when Handle is triggered", t, func() {
		t.SkipNow()
	})

	Convey("handler returns an error", t, func() {
		t.SkipNow()
	})
}

// TODO - complete unit test once SearchDataImportHandler.Handle has been implemented
func TestSearchDataImportHandler_Handle(t *testing.T) {
	Convey("Given a successful event handler, when Handle is triggered", t, func() {
		t.SkipNow()
	})

	Convey("handler returns an error", t, func() {
		t.SkipNow()
	})
}
