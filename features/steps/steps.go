package steps

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ONSdigital/dp-kafka/v3/kafkatest"
	"github.com/ONSdigital/dp-search-reindex-tracker/event"
	"github.com/ONSdigital/dp-search-reindex-tracker/service"
	"github.com/cucumber/godog"
	"github.com/rdumont/assistdog"
	"github.com/stretchr/testify/assert"
)

func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^these hello events are consumed:$`, c.theseHelloEventsAreConsumed)
	ctx.Step(`^I should receive a hello-world response$`, c.iShouldReceiveAHelloworldResponse)
}

func (c *Component) iShouldReceiveAHelloworldResponse() error {
	content, err := ioutil.ReadFile(c.cfg.OutputFilePath)
	if err != nil {
		return err
	}

	assert.Equal(c, "Hello, Tim!", string(content))

	return c.StepError()
}

func (c *Component) theseHelloEventsAreConsumed(table *godog.Table) error {

	observationEvents, err := c.convertToHelloEvents(table)
	if err != nil {
		return err
	}

	fmt.Println(observationEvents)

	signals := registerInterrupt()

	// run application in separate goroutine
	go func() {
		c.svc, err = service.Run(context.Background(), c.serviceList, "", "", "", c.errorChan)
	}()

	// consume extracted observations
	for _, e := range observationEvents {
		if err := c.sendToConsumer(e); err != nil {
			return err
		}
	}

	time.Sleep(300 * time.Millisecond)

	// kill application
	signals <- os.Interrupt

	return nil
}

func (c *Component) convertToHelloEvents(table *godog.Table) ([]*event.HelloCalledModel, error) {
	assist := assistdog.NewDefault()
	events, err := assist.CreateSlice(&event.HelloCalledModel{}, table)
	if err != nil {
		return nil, err
	}
	return events.([]*event.HelloCalledModel), nil
}

func (c *Component) sendToConsumer(e *event.HelloCalledModel) error {
	bytes, err := event.HelloCalledSchema.Marshal(e)
	if err != nil {
		return err
	}

	c.reindexTaskCountsConsumer.Channels().Upstream <- kafkatest.NewMessage(bytes, 0)
	return nil

}

func registerInterrupt() chan os.Signal {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	return signals
}
