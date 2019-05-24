Feature: List payments
  As an API user
  I want to get payments
  So I can get a list of payment transactions

  Background:
    Given I am authenticated to the API

  Scenario: List filtered payments
    When I send a "GET" request to "/v1/payments?filter[organisation_id]=743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb"
    Then the response code should be 200
    And the JSON response should contain payment data collection

  Scenario: List payments
    When I send a "GET" request to "/v1/payments"
    Then the response code should be 200
    And the JSON response should contain payment data collection
