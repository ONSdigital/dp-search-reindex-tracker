Feature: Consuming reindex-task-counts kafka messages

  Scenario: Posting and checking a response
    When these reindex-task-counts events are consumed:
        | JobID            | TaskName         | TaskCount      | TraceID          |
        | 1234             | test-component   | 2              | trace1234        |
    Then nothing happens