Feature: Fetch a payment
  As an API user
  I want to get a payment
  So I can get details of the payment transaction

  Background:
    Given I am authenticated to the API

  Scenario: Get payment
    When I send a "GET" request to "/v1/payments/a8dfdf10-33fa-4301-b859-e19853641655"
    And the JSON response should contain payment data
    Then the response code should be 200

  Scenario: Fail to get a payment
    When I send a "GET" request to "/v1/payments/aaaaaf10-33fa-4301-b859-e5555555555"
    Then the response code should be 404
    And the JSON response should contain an error
