Feature: Consuming reindex-task-counts kafka messages

  Scenario: Posting and checking a response
    When these reindex-task-counts events are consumed:
        | JobID            | TaskName         | ExtractionCompleted      | TaskCount  |
        | 1234             | test-component   | true                     | 2          |
    Then nothing happens