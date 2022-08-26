package event_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/ONSdigital/dp-search-reindex-tracker/config"
	"github.com/ONSdigital/dp-search-reindex-tracker/event/mock"

	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-kafka/v3/kafkatest"
	"github.com/ONSdigital/dp-search-reindex-tracker/event"
	. "github.com/smartystreets/goconvey/convey"
)

var testCtx = context.Background()

var errHandler = errors.New("handler error")

var testEvent = event.HelloCalledModel{
	RecipientName: "World",
}

// TODO: remove or replace hello called logic with app specific
func TestConsume(t *testing.T) {

	Convey("Given kafka consumer and event handler mocks", t, func() {

		cgChannels := &kafka.ConsumerGroupChannels{Upstream: make(chan kafka.Message, 2)}
		mockConsumer := &kafkatest.IConsumerGroupMock{
			ChannelsFunc: func() *kafka.ConsumerGroupChannels { return cgChannels },
		}

		handlerWg := &sync.WaitGroup{}
		mockEventHandler := &mock.HandlerMock{
			HandleFunc: func(ctx context.Context, config *config.Config, event *event.HelloCalledModel) error {
				defer handlerWg.Done()
				return nil
			},
		}

		mockEvent := &event.KafkaConsumerEvent[event.HelloCalledModel]{
			ConsumerGroup: mockConsumer,
			Handler:       mockEventHandler,
			Schema:        event.HelloCalledSchema,
		}

		Convey("And a kafka message with the valid schema being sent to the Upstream channel", func() {

			message := kafkatest.NewMessage(marshal(testEvent), 0)
			mockConsumer.Channels().Upstream <- message

			Convey("When consume message is called", func() {

				handlerWg.Add(1)
				event.Consume(testCtx, &config.Config{KafkaConfig: config.KafkaConfig{NumWorkers: 1}}, mockEvent)
				handlerWg.Wait()

				Convey("An event is sent to the mockEventHandler ", func() {
					So(len(mockEventHandler.HandleCalls()), ShouldEqual, 1)
					So(*mockEventHandler.HandleCalls()[0].HelloCalled, ShouldResemble, testEvent)
				})

				Convey("The message is committed and the consumer is released", func() {
					<-message.UpstreamDone()
					So(len(message.CommitCalls()), ShouldEqual, 1)
					So(len(message.ReleaseCalls()), ShouldEqual, 1)
				})
			})
		})

		Convey("And two kafka messages, one with a valid schema and one with an invalid schema", func() {

			validMessage := kafkatest.NewMessage(marshal(testEvent), 1)
			invalidMessage := kafkatest.NewMessage([]byte("invalid schema"), 0)
			mockConsumer.Channels().Upstream <- invalidMessage
			mockConsumer.Channels().Upstream <- validMessage

			Convey("When consume messages is called", func() {

				handlerWg.Add(1)
				event.Consume(testCtx, &config.Config{KafkaConfig: config.KafkaConfig{NumWorkers: 1}}, mockEvent)
				handlerWg.Wait()

				Convey("Only the valid event is sent to the mockEventHandler ", func() {
					So(len(mockEventHandler.HandleCalls()), ShouldEqual, 1)
					So(*mockEventHandler.HandleCalls()[0].HelloCalled, ShouldResemble, testEvent)
				})

				Convey("Only the valid message is committed, but the consumer is released for both messages", func() {
					<-validMessage.UpstreamDone()
					<-invalidMessage.UpstreamDone()
					So(len(validMessage.CommitCalls()), ShouldEqual, 1)
					So(len(invalidMessage.CommitCalls()), ShouldEqual, 1)
					So(len(validMessage.ReleaseCalls()), ShouldEqual, 1)
					So(len(invalidMessage.ReleaseCalls()), ShouldEqual, 1)
				})
			})
		})

		Convey("With a failing handler and a kafka message with the valid schema being sent to the Upstream channel", func() {
			mockEventHandler.HandleFunc = func(ctx context.Context, config *config.Config, event *event.HelloCalledModel) error {
				defer handlerWg.Done()
				return errHandler
			}
			mockEvent.Handler = mockEventHandler

			message := kafkatest.NewMessage(marshal(testEvent), 0)
			mockConsumer.Channels().Upstream <- message

			Convey("When consume message is called", func() {

				handlerWg.Add(1)
				event.Consume(testCtx, &config.Config{KafkaConfig: config.KafkaConfig{NumWorkers: 1}}, mockEvent)
				handlerWg.Wait()

				Convey("An event is sent to the mockEventHandler ", func() {
					So(len(mockEventHandler.HandleCalls()), ShouldEqual, 1)
					So(*mockEventHandler.HandleCalls()[0].HelloCalled, ShouldResemble, testEvent)
				})

				Convey("The message is committed and the consumer is released", func() {
					<-message.UpstreamDone()
					So(len(message.CommitCalls()), ShouldEqual, 1)
					So(len(message.ReleaseCalls()), ShouldEqual, 1)
				})
			})
		})
	})
}

// marshal helper method to marshal a event into a []byte
func marshal(e event.HelloCalledModel) []byte {
	bytes, err := event.HelloCalledSchema.Marshal(e)
	So(err, ShouldBeNil)
	return bytes
}
