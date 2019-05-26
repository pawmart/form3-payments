Feature: Create a payment
  As an API user
  I want to get payment
  So I can record a payment transaction

  Scenario: Create payment
    When I send a "POST" request to "/v1/payments" with:
    """
    {
        "data": {
            "organisation_id": "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb",
            "type": "Payment",
            "attributes": {
                "amount": "100.21",
                "currency": "GBP",
                "reference": "piano lessons",
                "beneficiary_party": {
                    "account_name": "W Owens",
                    "account_number": "31926819"
                },
                "debtor_party": {
                    "account_name": "EJ Brown Black",
                    "account_number": "GB29XABC10"
                }
            }
        }
    }
    """
    Then the response code should be 201
    And the JSON response should contain payment data

  Scenario: Fail to create payment
    When I send a "POST" request to "/v1/payments" with:
    """
        {
        }
    """
    Then the response code should be 400
    And the JSON response should contain an error
