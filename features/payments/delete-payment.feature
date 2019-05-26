Feature: Delete payment
  As an API user
  I want to delete a payment
  So I can remove a payment resource

  Scenario: Delete payment success
    When I send a "DELETE" request to "/v1/payments/b8dfdf10-33fa-4301-b859-e19853641651"
    Then the response code should be 204

  Scenario: Delete payment fail when not exist
    When I send a "DELETE" request to "/v1/payments/f3c5f34a-3985-44b2-bb1d-a51ffda32baf"
    Then the response code should be 404

  Scenario: Delete payment fail due to non uuid
    When I send a "DELETE" request to "/v1/payments/abc"
    Then the response code should be 422