package main

import (
	"context"
	"flag"
	"os"
	"testing"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-search-reindex-tracker/features/steps"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

var componentFlag = flag.Bool("component", false, "perform component tests")

type ComponentTest struct {
	MongoFeature *componenttest.MongoFeature
}

func (f *ComponentTest) InitializeScenario(godogCtx *godog.ScenarioContext) {
	ctx := context.Background()

	searchReindexTrackerComponent, err := steps.NewSearchReindexTrackerComponent()
	if err != nil {
		log.Error(ctx, "error occurred while creating a new search reindex tracker feature", err)
		os.Exit(1)
	}

	apiFeature := searchReindexTrackerComponent.InitAPIFeature()

	godogCtx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		searchReindexTrackerComponent.Reset()
		return ctx, nil
	})

	godogCtx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if err != nil {
			log.Error(ctx, "error retrieved after scenario", err)
			return ctx, err
		}

		searchReindexTrackerComponent.Close()

		return ctx, nil
	})

	searchReindexTrackerComponent.RegisterSteps(godogCtx)
	apiFeature.RegisterSteps(godogCtx)
}

func (f *ComponentTest) InitializeTestSuite(ctx *godog.TestSuiteContext) {
}

func TestComponent(t *testing.T) {
	if *componentFlag {
		status := 0

		var opts = godog.Options{
			Output: colors.Colored(os.Stdout),
			Format: "pretty",
			Paths:  flag.Args(),
		}

		f := &ComponentTest{}

		status = godog.TestSuite{
			Name:                 "feature_tests",
			ScenarioInitializer:  f.InitializeScenario,
			TestSuiteInitializer: f.InitializeTestSuite,
			Options:              &opts,
		}.Run()

		if status > 0 {
			t.Fail()
		}
	} else {
		t.Skip("component flag required to run component tests")
	}
}
