package main

// Payment definition.
type Payment struct {
	Type           string `json:"type";`
	ID             string `json:"id"`
	Version        int    `json:"version"`
	OrganisationID string `json:"organisation_id"`
	Attributes     *struct {
		Amount           string `json:"amount"`
		BeneficiaryParty *struct {
			AccountName   string `json:"account_name"`
			AccountNumber string `json:"account_number"`
		} `json:"beneficiary_party"`
		Currency    string `json:"currency"`
		DebtorParty *struct {
			AccountName   string `json:"account_name"`
			AccountNumber string `json:"account_number"`
		} `json:"debtor_party"`
		Reference string `json:"reference"`
	} `json:"attributes"`
}

// PaymentData definition.
type PaymentData struct {
	Data  *Payment `json:"data"`
	Links *Links   `json:"links"`
}

// PaymentDataList definition.
type PaymentDataList struct {
	Data  []*Payment `json:"data"`
	Links *Links     `json:"links"`
}

// APIError definition.
type APIError struct {
	Code    string `json:"error_code"`
	Message string `json:"error_message"`
}

// Links definition.
type Links struct {
	Self string `json:"self"`
}

//PaymentList definition.
type PaymentList struct {
	Data  []Payment `json:"data"`
	Links Links     `json:"links"`
}

// Health definition.
type Health struct {
	Status string `json:"status"`
}
