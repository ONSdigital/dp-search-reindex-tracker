package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	kafka "github.com/ONSdigital/dp-kafka/v3"
	"github.com/ONSdigital/dp-search-reindex-tracker/config"
	"github.com/ONSdigital/dp-search-reindex-tracker/event"
	"github.com/ONSdigital/log.go/v2/log"
)

const serviceName = "dp-search-reindex-tracker"

func main() {
	log.Namespace = serviceName
	ctx := context.Background()

	// Get Config
	cfg, err := config.Get()
	if err != nil {
		log.Fatal(ctx, "error getting config", err)
		os.Exit(1)
	}

	pConfig := &kafka.ProducerConfig{
		KafkaVersion: &cfg.KafkaConfig.Version,
	}
	if cfg.KafkaConfig.SecProtocol == config.KafkaTLSProtocolFlag {
		pConfig.SecurityConfig = kafka.GetSecurityConfig(
			cfg.KafkaConfig.SecCACerts,
			cfg.KafkaConfig.SecClientCert,
			cfg.KafkaConfig.SecClientKey,
			cfg.KafkaConfig.SecSkipVerify,
		)
	}
	kafkaProducer, err := kafka.NewProducer(ctx, pConfig)
	if err != nil {
		log.Fatal(ctx, "fatal error trying to create kafka producer", err, log.Data{"topic": cfg.KafkaConfig.ReindexRequestedTopic})
		os.Exit(1)
	}

	// kafka error logging go-routines
	kafkaProducer.LogErrors(ctx)

	time.Sleep(500 * time.Millisecond)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		e := scanEvent(scanner)
		log.Info(ctx, "sending reindex-requested event", log.Data{"reindexRequestedEvent": e})

		bytes, err := event.ReindexRequestedSchema.Marshal(e)
		if err != nil {
			log.Fatal(ctx, "reindex-requested event error", err)
			os.Exit(1)
		}

		// Send bytes to Output channel, after calling Initialise just in case it is not initialised.
		kafkaProducer.Initialise(ctx)
		kafkaProducer.Channels().Output <- bytes
	}
}

// scanEvent creates a ReindexRequested event according to the user input
func scanEvent(scanner *bufio.Scanner) *event.ReindexRequestedModel {
	fmt.Println("--- [Send Kafka ReindexRequested] ---")

	fmt.Println("Please type the job id")
	fmt.Printf("$ ")
	scanner.Scan()
	jobID := scanner.Text()

	return &event.ReindexRequestedModel{
		JobID:       jobID,
		SearchIndex: "test",
		TraceID:     "trace1234",
	}
}
