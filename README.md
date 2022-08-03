# dp-search-reindex-tracker
dp-search-reindex-tracker tracks a search reindex job

### Getting started

* Run `make debug`

The Search Reindex Tracker service runs in the background consuming messages from Kafka. 

The following diagram shows how the Search Reindex Tracker will be used to track the progress of search reindex jobs and their tasks:
[Search Reindex Pipeline](https://miro.com/app/board/o9J_liKF4Lg=/?fromRedirect=1)

The topics will produce the following messages:
- The Search Reindex API will produce messages, for the reindex-requested topic, that request new search reindex jobs. 
  The Search Reindex tracker will consume these and set the state of the relevant search reindex jobs to "in-progress". 
- The Search Data Finder will produce messages for the reindex-task-counts topic. 
  The Search Reindex tracker will consume these in order to set the "total search documents" number for the relevant search reindex jobs. 
- The Search Data Importer will produce messages for the search-data-imported topic.
  The Search Reindex Tracker will consume these in order to set the "total inserted search documents" number for the relevant search reindex jobs.

### Dependencies

* Requires running…
  * [kafka](https://github.com/ONSdigital/dp/blob/main/guides/INSTALLING.md#prerequisites)
* No further dependencies other than those defined in `go.mod`

### Configuration

| Environment variable               | Default                           | Description
| ----------------------------       | --------------------------------- | -----------
| BIND_ADDR                          | localhost:28500                    | The host and port to bind to
| GRACEFUL_SHUTDOWN_TIMEOUT          | 5s                                | The graceful shutdown timeout in seconds (`time.Duration` format)
| HEALTHCHECK_INTERVAL               | 30s                               | Time between self-healthchecks (`time.Duration` format)
| HEALTHCHECK_CRITICAL_TIMEOUT       | 90s                               | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format)
| KAFKA_ADDR                         | "localhost:9092"                  | The address of Kafka (accepts list)
| KAFKA_OFFSET_OLDEST                | true                              | Start processing Kafka messages in order from the oldest in the queue
| KAFKA_NUM_WORKERS                  | 1                                 | The maximum number of parallel kafka consumers
| KAFKA_REINDEX_REQUESTED_GROUP      | dp-search-reindex-tracker         | The consumer group name for reindex requested events
| KAFKA_REINDEX_REQUESTED_TOPIC      | reindex-requested                 | The topic name for reindex requested events
| KAFKA_REINDEX_TASK_COUNTS_GROUP    | dp-search-reindex-tracker         | The consumer group name for reindex task count events
| KAFKA_REINDEX_TASK_COUNTS_TOPIC    | reindex-task-counts               | The topic name for reindex task count events
| KAFKA_SEARCH_DATA_IMPORTED_GROUP   | dp-search-reindex-tracker         | The consumer group name for search data imported events
| KAFKA_SEARCH_DATA_IMPORTED_TOPIC   | search-data-imported              | The topic name for search data imported events
| KAFKA_SEC_PROTO                    | _unset_                           | if set to `TLS`, kafka connections will use TLS ([kafka TLS doc])
| KAFKA_SEC_CA_CERTS                 | _unset_                           | CA cert chain for the server cert ([kafka TLS doc])
| KAFKA_SEC_CLIENT_KEY               | _unset_                           | PEM for the client key ([kafka TLS doc])
| KAFKA_SEC_CLIENT_CERT              | _unset_                           | PEM for the client certificate ([kafka TLS doc])
| KAFKA_SEC_SKIP_VERIFY              | false                             | ignores server certificate issues if `true` ([kafka TLS doc])

[kafka TLS doc]: https://github.com/ONSdigital/dp-kafka/tree/main/examples#tls

### Healthcheck

 The `/health` endpoint returns the current status of the service. Dependent services are health checked on an interval defined by the `HEALTHCHECK_INTERVAL` environment variable.

 On a development machine a request to the health check endpoint can be made by:

 `curl localhost:8125/health`

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2022, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
