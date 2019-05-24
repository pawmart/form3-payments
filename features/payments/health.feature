Feature: Check health
  As an API user
  I want to check API health
  So I can be sure I can interact with the API

  Background:
    Given I am authenticated to the API
    When I send a "GET" request to "/health"

  Scenario: Success
    Then the response code should be 200
    And the response should have health status up