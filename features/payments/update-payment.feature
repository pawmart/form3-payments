Feature: Update payment
  As an API user
  I want to update a payment
  So I can change the state of the payment

  Background:
    Given I am authenticated to the API


  Scenario: Update payment
    When I send a "PATCH" request to "/v1/payments" with:
    """
    {
        "data": {
            "id": "a8dfdf10-33fa-4301-b859-e19853641655",
            "organisation_id": "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb",
            "type": "Payment",
            "attributes": {
                "reference": "abc"
            }
        }
    }
    """
    And the response code should be 200