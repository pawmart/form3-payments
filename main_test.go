package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/pawmart/form3-payments/models"
	"github.com/pawmart/form3-payments/restapi/operations"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/pawmart/form3-payments/storage"
)

var handler http.Handler

type featureState struct {
	api    *operations.Form3paymentsAPI
	resp   *httptest.ResponseRecorder
	apikey string
}

func (a *featureState) resetState(interface{}) {
	a.resp = httptest.NewRecorder()
	a.apikey = ""
}

func (a *featureState) prepareSuite() {

	handler = getAPIServer().GetHandler()

	db, _ := storage.New().Db()
	db.DropDatabase()

	recordJSON := `{
            "id" : "b8dfdf10-33fa-4301-b859-e19853641651",
            "type": "Payments",
            "organisation_id": "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb",
            "created_on": 1558772378,
            "modified_on": 1558772378,
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
        }`

	recordJSON2 := `{
            "id" : "a8dfdf10-33fa-4301-b859-e19853641655",
            "type": "Payments",
            "organisation_id": "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb",
            "created_on": 1558772378,
            "modified_on": 1558772378,
            "attributes": {
                "amount": "3.21",
                "currency": "GBP",
                "reference": "boxing lessons",
                "beneficiary_party": {
                    "account_name": "W Doe",
                    "account_number": "31926819"
                },
                "debtor_party": {
                    "account_name": "PJ Yo Black",
                    "account_number": "GB29XABC10"
                }
            }
        }`

	p := new(Payment)
	p2 := new(Payment)
	json.Unmarshal([]byte(recordJSON), p)
	json.Unmarshal([]byte(recordJSON2), p2)

	err := db.C("payments").Insert(p)
	if err != nil {
		fmt.Errorf("populate db failed")
	}
	err = db.C("payments").Insert(p2)
	if err != nil {
		fmt.Errorf("populate db failed")
	}
}

func (a *featureState) callEndpoint(method string, endpoint string, req *http.Request) {

	if a.apikey != "" {
		req.Header.Set("X-APIKEY", "abc")
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	a.resp = rr
}

func (a *featureState) iSendRequestTo(method, endpoint string) (err error) {
	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	a.callEndpoint(method, endpoint, req)
	return
}

func (a *featureState) iSendARequestToWith(method, endpoint string, body *gherkin.DocString) (err error) {
	req, err := http.NewRequest(method, endpoint, strings.NewReader(body.Content))
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	if a.apikey != "" {
		req.Header.Set("X-APIKEY", "abc")
	}
	a.callEndpoint(method, endpoint, req)
	return
}

func (a *featureState) theResponseCodeShouldBe(code int) error {
	if code != a.resp.Code {
		return fmt.Errorf("expected response code to be: %d, but actual is: %d", code, a.resp.Code)
	}
	return nil
}

func (a *featureState) theJSONResponseShouldContainAnError() (err error) {
	var error models.APIError
	if err = json.Unmarshal([]byte(a.resp.Body.String()), &error); err != nil {
		return fmt.Errorf(err.Error())
	}
	if err != nil {
		return fmt.Errorf("could not unmarshal error response", error, err.Error(), a.resp)
	}
	if len(error.ErrorMessage) > 0 && len(error.ErrorCode) > 0 {
		return
	} else {
		return fmt.Errorf("no error message")
	}
	return
}

func (a *featureState) theJSONResponseShouldContainPaymentData() (err error) {
	var pd *models.PaymentCreationResponse
	if err = json.Unmarshal([]byte(a.resp.Body.String()), &pd); err != nil {
		return fmt.Errorf(err.Error())
	}
	if err != nil {
		return fmt.Errorf("could not unmarshal payment data response", pd, err.Error(), a.resp.Body.String())
	}

	log.Print(pd)

	if pd.Data != nil && len(pd.Data.ID) > 0 {
		return
	}
	return
}

func (a *featureState) theResponseShouldHaveHealthStatusUp() (err error) {
	var pd *models.Health
	if err = json.Unmarshal([]byte(a.resp.Body.String()), &pd); err != nil {
		return fmt.Errorf(err.Error())
	}

	if err != nil {
		return fmt.Errorf("could not unmarshal response", pd, err.Error(), a.resp.Body.String())
	}
	if pd.Status != "up" {
		return fmt.Errorf("health status is not up", pd, err.Error(), a.resp.Body.String())
	}
	return
}

func (a *featureState) theJSONResponseShouldContainPaymentDataWithReferenceAttributeOf(arg1 string) (err error) {
	var pd *models.PaymentDetailsResponse
	if err = json.Unmarshal([]byte(a.resp.Body.String()), &pd); err != nil {
		return fmt.Errorf(err.Error())
	}
	if err != nil {
		return fmt.Errorf("could not unmarshal response", pd, err.Error(), a.resp.Body.String())
	}
	if pd.Data.Attributes.Reference != arg1 {
		return fmt.Errorf("wrong reference", pd, err.Error(), a.resp.Body.String())
	}
	return
}

func (a *featureState) theJSONResponseShouldContainPaymentDataCollection() (err error) {
	var pd *models.PaymentDetailsListResponse
	if err = json.Unmarshal([]byte(a.resp.Body.String()), &pd); err != nil {
		return fmt.Errorf(err.Error())
	}

	if pd.Data == nil || len(pd.Data) < 1 {
		return fmt.Errorf("no collection returned")
	}
	if pd.Data[0].ID == "" {
		return fmt.Errorf("no collection record returned")
	}

	return
}

func (a *featureState) theJSONResponseShouldContainNoPaymentDataCollection() (err error) {
	var pd *models.PaymentDetailsListResponse
	if err = json.Unmarshal([]byte(a.resp.Body.String()), &pd); err != nil {
		fmt.Errorf(err.Error())
	}

	log.Print(pd)
	log.Print(pd.Data)

	for _, v := range pd.Data {
		return fmt.Errorf("payment data exists and it should not", pd, a.resp.Body.String(), v)
	}

	return
}

func (a *featureState) iAmAuthenticatedToTheAPI() error {

	a.apikey = "abc"
	return nil
}

func FeatureContext(s *godog.Suite) {
	api := &featureState{}

	s.BeforeSuite(api.prepareSuite)
	s.BeforeScenario(api.resetState)

	s.Step(`^I am authenticated to the API$`, api.iAmAuthenticatedToTheAPI)
	s.Step(`^I send a "([^"]*)" request to "([^"]*)"$`, api.iSendRequestTo)
	s.Step(`^I send a "([^"]*)" request to "([^"]*)" with:$`, api.iSendARequestToWith)
	s.Step(`^the response code should be (\d+)$`, api.theResponseCodeShouldBe)
	s.Step(`^the JSON response should contain an error$`, api.theJSONResponseShouldContainAnError)
	s.Step(`^the JSON response should contain payment data$`, api.theJSONResponseShouldContainPaymentData)
	s.Step(`^the JSON response should contain no payment data collection$`, api.theJSONResponseShouldContainNoPaymentDataCollection)
	s.Step(`^the response should have health status up$`, api.theResponseShouldHaveHealthStatusUp)
	s.Step(`^the JSON response should contain payment data collection$`, api.theJSONResponseShouldContainPaymentDataCollection)
	s.Step(`^the JSON response should contain payment data with reference attribute of "([^"]*)"$`, api.theJSONResponseShouldContainPaymentDataWithReferenceAttributeOf)
}
