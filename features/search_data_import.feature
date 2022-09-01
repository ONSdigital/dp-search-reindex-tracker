Feature: Consuming search-data-import kafka messages

  Scenario: Posting and checking a response
    When these search-data-import events are consumed:
      | UID              | Title            | TraceID          |
      | 1234             | test-component   | trace1234        |
    Then nothing happens