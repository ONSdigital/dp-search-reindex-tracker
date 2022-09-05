package steps

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/ONSdigital/dp-search-reindex-tracker/event"
	"github.com/cucumber/godog"
)

func (c *SearchReindexTrackerComponent) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^all of the downstream services are healthy$`, c.allOfTheDownstreamServicesAreHealthy)
	ctx.Step(`^I should receive the following health JSON response:$`, c.iShouldReceiveTheFollowingHealthJSONResponse)
	ctx.Step(`^I wait (\d+) seconds`, c.delayTimeBySeconds)
	ctx.Step(`^nothing happens`, c.nothingHappens)
	ctx.Step(`^one of the downstream services is failing`, c.oneOfTheDownstreamServicesIsFailing)
	ctx.Step(`^one of the downstream services is warning`, c.oneOfTheDownstreamServicesIsWarning)
	ctx.Step(`^these reindex-requested events are consumed:$`, c.theseReindexRequestedEventsAreConsumed)
	ctx.Step(`^these reindex-task-counts events are consumed:$`, c.theseReindexTaskCountsEventsAreConsumed)
	ctx.Step(`^these search-data-import events are consumed:$`, c.theseSearchDataImportEventsAreConsumed)
}

func (c *SearchReindexTrackerComponent) allOfTheDownstreamServicesAreHealthy() error {
	c.fakeAPIRouter.healthRequest.Lock()
	defer c.fakeAPIRouter.healthRequest.Unlock()

	c.fakeAPIRouter.healthRequest.CustomHandle = healthCheckStatusHandle(200)

	return nil
}

// delayTimeBySeconds pauses the goroutine for the given seconds
func (c *SearchReindexTrackerComponent) delayTimeBySeconds(sec int) error {
	time.Sleep(time.Duration(int64(sec)) * time.Second)
	return nil
}

func (c *SearchReindexTrackerComponent) iShouldReceiveTheFollowingHealthJSONResponse(expectedResponse *godog.DocString) error {
	var healthResponse, expectedHealth HealthCheckTest

	responseBody, err := ioutil.ReadAll(c.apiFeature.HttpResponse.Body)
	if err != nil {
		return fmt.Errorf("failed to read response of search controller component - error: %v", err)
	}

	err = json.Unmarshal(responseBody, &healthResponse)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response of search controller component - error: %v", err)
	}

	err = json.Unmarshal([]byte(expectedResponse.Content), &expectedHealth)
	if err != nil {
		return fmt.Errorf("failed to unmarshal expected health response - error: %v", err)
	}

	c.validateHealthCheckResponse(healthResponse, expectedHealth)

	return c.ErrorFeature.StepError()
}

func (c *SearchReindexTrackerComponent) nothingHappens() error {
	return nil
}

func (c *SearchReindexTrackerComponent) oneOfTheDownstreamServicesIsWarning() error {
	c.fakeAPIRouter.healthRequest.Lock()
	defer c.fakeAPIRouter.healthRequest.Unlock()

	c.fakeAPIRouter.healthRequest.CustomHandle = healthCheckStatusHandle(429)

	return nil
}

func (c *SearchReindexTrackerComponent) oneOfTheDownstreamServicesIsFailing() error {
	c.fakeAPIRouter.healthRequest.Lock()
	defer c.fakeAPIRouter.healthRequest.Unlock()

	c.fakeAPIRouter.healthRequest.CustomHandle = healthCheckStatusHandle(500)

	return nil
}

func (c *SearchReindexTrackerComponent) theseReindexRequestedEventsAreConsumed(table *godog.Table) error {
	return theseEventsAreConsumed[event.ReindexRequestedModel](c, event.ReindexRequestedSchema, c.reindexRequestedConsumer, table)
}

func (c *SearchReindexTrackerComponent) theseReindexTaskCountsEventsAreConsumed(table *godog.Table) error {
	return theseEventsAreConsumed[event.ReindexTaskCountsModel](c, event.ReindexTaskCountsSchema, c.reindexTaskCountsConsumer, table)
}

func (c *SearchReindexTrackerComponent) theseSearchDataImportEventsAreConsumed(table *godog.Table) error {
	return theseEventsAreConsumed[event.SearchDataImportModel](c, event.SearchDataImportSchema, c.searchDataImportedConsumer, table)
}
