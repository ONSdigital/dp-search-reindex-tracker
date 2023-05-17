package event_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ONSdigital/dp-kafka/v3/kafkatest"
	"github.com/ONSdigital/dp-search-reindex-tracker/config"
	"github.com/ONSdigital/dp-search-reindex-tracker/event"
	"github.com/ONSdigital/dp-search-reindex-tracker/event/mock"
	"github.com/ONSdigital/go-ns/avro"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	testCtx = context.Background()

	testReindexRequestedEvent = event.ReindexRequestedModel{
		JobID:       "1234",
		SearchIndex: "test",
		TraceID:     "5678",
	}

	handlerSuccess = &mock.HandlerMock{
		HandleFunc: func(ctx context.Context, cfg *config.Config, reindexRequested *event.ReindexRequestedModel) error {
			return nil
		},
	}

	errHandler = errors.New("handler error")

	handlerFail = &mock.HandlerMock{
		HandleFunc: func(ctx context.Context, cfg *config.Config, reindexRequested *event.ReindexRequestedModel) error {
			return errHandler
		},
	}
)

func TestProcessMessage(t *testing.T) {
	cfg, err := config.Get()
	if err != nil {
		t.Errorf("failed to get config for testing")
	}

	Convey("Given a valid message", t, func() {
		topicEvent := &event.KafkaConsumerEvent[event.ReindexRequestedModel]{
			Handler: handlerSuccess,
			Schema:  event.ReindexRequestedSchema,
		}

		validMsg := marshal(event.ReindexRequestedSchema, testReindexRequestedEvent)
		kafkaMsg, err := kafkatest.NewMessage(validMsg, 0)
		So(err, ShouldBeNil)

		Convey("When ProcessMessage is called", func() {
			err := event.ProcessMessage(testCtx, cfg, topicEvent, kafkaMsg)

			Convey("Then message is processed and no error is returned", func(c C) {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given an invalid message", t, func() {
		topicEvent := &event.KafkaConsumerEvent[event.ReindexRequestedModel]{
			Handler: handlerSuccess,
			Schema:  event.ReindexRequestedSchema,
		}

		kafkaMsg, err := kafkatest.NewMessage([]byte("invalid"), 0)
		So(err, ShouldBeNil)

		Convey("When ProcessMessage is called", func() {
			err := event.ProcessMessage(testCtx, cfg, topicEvent, kafkaMsg)

			Convey("Then an error is returned", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("Given the topicEvent is nil", t, func() {
		var topicEvent *event.KafkaConsumerEvent[event.ReindexRequestedModel]

		validMsg := marshal(event.ReindexRequestedSchema, testReindexRequestedEvent)
		kafkaMsg, err := kafkatest.NewMessage(validMsg, 0)
		So(err, ShouldBeNil)

		Convey("When ProcessMessage is called", func() {
			err := event.ProcessMessage(testCtx, cfg, topicEvent, kafkaMsg)

			Convey("Then an error is returned", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("Given the handler returns an error", t, func() {
		topicEvent := &event.KafkaConsumerEvent[event.ReindexRequestedModel]{
			Handler: handlerFail,
			Schema:  event.ReindexRequestedSchema,
		}

		validMsg := marshal(event.ReindexRequestedSchema, testReindexRequestedEvent)
		kafkaMsg, err := kafkatest.NewMessage(validMsg, 0)
		So(err, ShouldBeNil)

		Convey("When ProcessMessage is called", func() {
			err := event.ProcessMessage(testCtx, cfg, topicEvent, kafkaMsg)

			Convey("Then an error is returned", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

// marshal helper method to marshal a event into a []byte
func marshal[M event.KafkaAvroModel](topicSchema *avro.Schema, data M) []byte {
	bytes, err := topicSchema.Marshal(data)
	So(err, ShouldBeNil)
	return bytes
}
