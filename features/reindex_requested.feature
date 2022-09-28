Feature: Consuming reindex-requested kafka messages 

  Scenario: Patch request to search-reindex-api is successful
    Given all of the downstream services are healthy
    And patch request to search-reindex-api is successful for job id "1234"
    When these reindex-requested events are consumed:
      | JobID            | SearchIndex      | TraceID          |
      | 1234             | test-component   | trace1234        |
    Then the state of the reindex job should be updated to in-progress

  Scenario: Error received from patch request to search-reindex-api
    Given all of the downstream services are healthy
    And patch request to search-reindex-api is unsuccessful for job id "1234"
    When these reindex-requested events are consumed:
      | JobID            | SearchIndex      | TraceID          |
      | 1234             | test-component   | trace1234        |
    Then the state of the reindex job should not be updated to in-progress