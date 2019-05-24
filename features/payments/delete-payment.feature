Feature: Delete payment
  As an API user
  I want to delete a payment
  So I can remove a payment resource

  Background:
    Given I am authenticated to the API

  Scenario: Delete payment success
    When I send a "DELETE" request to "/v1/payments/b8dfdf10-33fa-4301-b859-e19853641651"
    Then the response code should be 204
