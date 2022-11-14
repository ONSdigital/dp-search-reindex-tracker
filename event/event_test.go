package event

import (
	"fmt"
	"testing"

	"github.com/ONSdigital/dp-kafka/v3/kafkatest"
	searchReindexAPIMock "github.com/ONSdigital/dp-search-reindex-api/sdk/mocks"
	. "github.com/smartystreets/goconvey/convey"
)

func TestValidate(t *testing.T) {
	Convey("Given valid options for kafka event", t, func() {
		options := &KafkaEventOptions{
			ConsumerGroup:          &kafkatest.IConsumerGroupMock{},
			SearchReindexAPIClient: &searchReindexAPIMock.ClientMock{},
		}

		Convey("When validate is called", func() {
			err := options.validate()

			Convey("Then no error should be returned", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given options is nil for kafka event", t, func() {
		var options *KafkaEventOptions

		Convey("When validate is called", func() {
			err := options.validate()

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, fmt.Errorf("options for event not provided"))
			})
		})
	})

	Convey("Given a nil consumer group", t, func() {
		options := &KafkaEventOptions{
			SearchReindexAPIClient: &searchReindexAPIMock.ClientMock{},
		}

		Convey("When validate is called", func() {
			err := options.validate()

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, fmt.Errorf("consumer group not provided"))
			})
		})
	})
}

func TestGetReindexRequested(t *testing.T) {
	Convey("Given valid options for reindex requested event", t, func() {
		options := &KafkaEventOptions{
			ConsumerGroup:          &kafkatest.IConsumerGroupMock{},
			SearchReindexAPIClient: &searchReindexAPIMock.ClientMock{},
		}

		Convey("When GetReindexRequested is called", func() {
			event, err := GetReindexRequested(options)

			Convey("Then information needed for reindex requested event is returned", func() {
				So(event, ShouldNotBeNil)
				So(event, ShouldNotBeEmpty)

				So(event.Handler, ShouldResemble, &ReindexRequestedHandler{
					SearchReindexAPIClient: options.SearchReindexAPIClient,
				})
				So(event.Schema, ShouldResemble, ReindexRequestedSchema)
				So(event.ConsumerGroup, ShouldResemble, options.ConsumerGroup)

				Convey("And no error should be returned", func() {
					So(err, ShouldBeNil)
				})
			})
		})
	})

	Convey("Given invalid options for reindex requested event", t, func() {
		// consumer group is nil
		options := &KafkaEventOptions{
			SearchReindexAPIClient: &searchReindexAPIMock.ClientMock{},
		}

		Convey("When GetReindexRequested is called", func() {
			event, err := GetReindexRequested(options)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, fmt.Errorf("failed to validate options given for reindex requested event"))

				Convey("And object returned should be nil", func() {
					So(event, ShouldBeNil)
				})
			})
		})
	})

	Convey("Given search reindex api client is nil", t, func() {
		options := &KafkaEventOptions{
			ConsumerGroup: &kafkatest.IConsumerGroupMock{},
		}

		Convey("When GetReindexRequested is called", func() {
			event, err := GetReindexRequested(options)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, fmt.Errorf("search reindex api client not provided"))

				Convey("And object returned should be nil", func() {
					So(event, ShouldBeNil)
				})
			})
		})
	})
}

func TestGetReindexTaskCounts(t *testing.T) {
	Convey("Given valid options for reindex task counts event", t, func() {
		options := &KafkaEventOptions{
			ConsumerGroup:          &kafkatest.IConsumerGroupMock{},
			SearchReindexAPIClient: &searchReindexAPIMock.ClientMock{},
		}

		Convey("When GetReindexTaskCounts is called", func() {
			event, err := GetReindexTaskCounts(options)

			Convey("Then information needed for reindex task counts event is returned", func() {
				So(event, ShouldNotBeNil)
				So(event, ShouldNotBeEmpty)

				So(event.Handler, ShouldResemble, &ReindexTaskCountsHandler{
					SearchReindexAPIClient: options.SearchReindexAPIClient,
				})
				So(event.Schema, ShouldResemble, ReindexTaskCountsSchema)
				So(event.ConsumerGroup, ShouldResemble, options.ConsumerGroup)

				Convey("And no error should be returned", func() {
					So(err, ShouldBeNil)
				})
			})
		})
	})

	Convey("Given invalid options for reindex task counts event", t, func() {
		// consumer group is nil
		options := &KafkaEventOptions{}

		Convey("When GetReindexTaskCounts is called", func() {
			event, err := GetReindexTaskCounts(options)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, fmt.Errorf("failed to validate options given for reindex task counts event"))

				Convey("And object returned should be nil", func() {
					So(event, ShouldBeNil)
				})
			})
		})
	})
}

func TestGetSearchDataImport(t *testing.T) {
	Convey("Given valid options for search data import event", t, func() {
		options := &KafkaEventOptions{
			ConsumerGroup: &kafkatest.IConsumerGroupMock{},
		}

		Convey("When GetSearchDataImport is called", func() {
			event, err := GetSearchDataImport(options)

			Convey("Then information needed for search data import event is returned", func() {
				So(event, ShouldNotBeNil)
				So(event, ShouldNotBeEmpty)

				So(event.Handler, ShouldResemble, &SearchDataImportHandler{})
				So(event.Schema, ShouldResemble, SearchDataImportSchema)
				So(event.ConsumerGroup, ShouldResemble, options.ConsumerGroup)

				Convey("And no error should be returned", func() {
					So(err, ShouldBeNil)
				})
			})
		})
	})

	Convey("Given invalid options for search data import event", t, func() {
		// consumer group is nil
		options := &KafkaEventOptions{}

		Convey("When GetSearchDataImport is called", func() {
			event, err := GetSearchDataImport(options)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, fmt.Errorf("failed to validate options given for search data import event"))

				Convey("And object returned should be nil", func() {
					So(event, ShouldBeNil)
				})
			})
		})
	})
}
