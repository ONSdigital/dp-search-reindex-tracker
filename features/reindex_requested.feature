Feature: Consuming reindex-requested kafka messages 

  Scenario: Posting and checking a response
    When these reindex-requested events are consumed:
      | JobID            | SearchIndex      | TraceID          |
      | 1234             | test-component   | trace1234        |
    Then nothing happens