package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/gin-gonic/gin"
	"github.com/pawmart/form3-payments/storage"
)

type featureState struct {
	resp   *httptest.ResponseRecorder
	router *gin.Engine
	apikey string
}

func (a *featureState) resetState(interface{}) {
	a.resp = httptest.NewRecorder()
	a.apikey = ""
}

func (a *featureState) prepareSuite() {

	a.router = setupRouter()

	db, _ := storage.New().Db()
	db.DropDatabase()

	recordJSON := `{
            "id" : "b8dfdf10-33fa-4301-b859-e19853641651",
            "type": "Payments",
            "organisation_id": "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb",
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

func (a *featureState) iSendRequestTo(method, endpoint string) (err error) {
	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return
	}
	if a.apikey != "" {
		req.Header.Set("X-APIKEY", "abc")
	}
	a.router.ServeHTTP(a.resp, req)
	return
}

func (a *featureState) iSendARequestToWith(method, endpoint string, body *gherkin.DocString) (err error) {
	req, err := http.NewRequest(method, endpoint, strings.NewReader(body.Content))
	if err != nil {
		return
	}
	if a.apikey != "" {
		req.Header.Set("X-APIKEY", "abc")
	}
	a.router.ServeHTTP(a.resp, req)
	return
}

func (a *featureState) theResponseCodeShouldBe(code int) error {
	if code != a.resp.Code {
		return fmt.Errorf("expected response code to be: %d, but actual is: %d", code, a.resp.Code)
	}
	return nil
}

func (a *featureState) theJSONResponseShouldContainAnError() (err error) {
	var error APIError
	if err = json.Unmarshal([]byte(a.resp.Body.String()), &error); err != nil {
		return
	}
	if err != nil {
		return fmt.Errorf("could not unmarshal error response", error, err.Error(), a.resp.Body.String())
	}
	if len(error.Message) > 0 && len(error.Code) > 0 {
		return
	}
	return
}

func (a *featureState) theJSONResponseShouldContainPaymentData() (err error) {
	var pd *PaymentData
	if err = json.Unmarshal([]byte(a.resp.Body.String()), &pd); err != nil {
		return
	}
	if err != nil {
		return fmt.Errorf("could not unmarshal payment data response", pd, err.Error(), a.resp.Body.String())
	}
	if pd.Data != nil && len(pd.Data.ID) > 0 {
		return
	}
	return
}

func (a *featureState) theResponseShouldHaveHealthStatusUp() (err error) {
	var pd *Health
	if err = json.Unmarshal([]byte(a.resp.Body.String()), &pd); err != nil {
		return
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
	var pd *PaymentData
	if err = json.Unmarshal([]byte(a.resp.Body.String()), &pd); err != nil {
		return
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
	var pd *PaymentDataList
	if err = json.Unmarshal([]byte(a.resp.Body.String()), &pd); err != nil {
		return
	}

	if pd.Data[0].ID == "" {
		return fmt.Errorf("no collection record returned")
	}
	if len(pd.Data) < 2 {
		return fmt.Errorf("no collection returned")
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
	s.Step(`^the response should have health status up$`, api.theResponseShouldHaveHealthStatusUp)
	s.Step(`^the JSON response should contain payment data collection$`, api.theJSONResponseShouldContainPaymentDataCollection)
	s.Step(`^the JSON response should contain payment data with reference attribute of "([^"]*)"$`, api.theJSONResponseShouldContainPaymentDataWithReferenceAttributeOf)
}
