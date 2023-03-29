package event

import (
	"context"
	"errors"
	"testing"

	"github.com/ONSdigital/dp-search-reindex-api/sdk"
	searchReindexAPIClientMock "github.com/ONSdigital/dp-search-reindex-api/sdk/mocks"
	"github.com/ONSdigital/dp-search-reindex-tracker/config"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	testReindexRequestEventOK = &ReindexRequestedModel{
		JobID:       "1234",
		SearchIndex: "testing",
		TraceID:     "5678",
	}
	testReindexTaskCountEventOK = &ReindexTaskCountsModel{
		JobID:   "1234",
		TraceID: "5678",
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
		searchReindexClientMock := &searchReindexAPIClientMock.ClientMock{
			PatchJobFunc: func(ctx context.Context, headers sdk.Headers, jobID string, body []sdk.PatchOperation) (*sdk.RespHeaders, error) {
				return nil, nil
			},
		}

		reindexRequestedHandler := &ReindexRequestedHandler{
			SearchReindexAPIClient: searchReindexClientMock,
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
		searchReindexClientMock := &searchReindexAPIClientMock.ClientMock{
			PatchJobFunc: func(ctx context.Context, headers sdk.Headers, jobID string, body []sdk.PatchOperation) (*sdk.RespHeaders, error) {
				return nil, errUnexpected
			},
		}

		reindexRequestedHandler := &ReindexRequestedHandler{
			SearchReindexAPIClient: searchReindexClientMock,
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

func TestReindexTaskCountsHandler_Handle(t *testing.T) {
	ctx := context.Background()

	cfg, err := config.Get()
	if err != nil {
		t.Errorf("failed to get config - err: %v", err)
	}

	Convey("Given post request to search reindex api is not successful", t, func() {
		searchReindexClientMock := &searchReindexAPIClientMock.ClientMock{
			PatchJobFunc: func(ctx context.Context, headers sdk.Headers, jobID string, body []sdk.PatchOperation) (*sdk.RespHeaders, error) {
				return nil, nil
			},
			PutTaskNumberOfDocsFunc: func(ctx context.Context, reqHeaders sdk.Headers, jobID string, taskName string, docCount string) (*sdk.RespHeaders, error) {
				return nil, errUnexpected
			},
		}

		reindexTaskCountHandler := &ReindexTaskCountsHandler{
			SearchReindexAPIClient: searchReindexClientMock,
		}

		Convey("When ReindexTaskCountHandler.Handle is called", func() {
			err := reindexTaskCountHandler.Handle(ctx, cfg, testReindexTaskCountEventOK)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given request to search reindex api is successful", t, func() {
		searchReindexClientMock := &searchReindexAPIClientMock.ClientMock{
			PatchJobFunc: func(ctx context.Context, headers sdk.Headers, jobID string, body []sdk.PatchOperation) (*sdk.RespHeaders, error) {
				return nil, nil
			},
			PutTaskNumberOfDocsFunc: func(ctx context.Context, reqHeaders sdk.Headers, jobID string, taskName string, docCount string) (*sdk.RespHeaders, error) {
				return nil, nil
			},
		}

		reindexTaskCountHandler := &ReindexTaskCountsHandler{
			SearchReindexAPIClient: searchReindexClientMock,
		}

		Convey("When ReindexTaskCountHandler.Handle is called", func() {
			err := reindexTaskCountHandler.Handle(ctx, cfg, testReindexTaskCountEventOK)
			So(err, ShouldBeNil)
		})
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
