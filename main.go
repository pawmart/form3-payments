package main

import (
	"log"

	"github.com/alexflint/go-arg"
	"github.com/go-openapi/loads"

	"github.com/pawmart/form3-payments/restapi"
	"github.com/pawmart/form3-payments/restapi/operations"
)

type cliArgs struct {
	Port int `arg:"-p,help:port to listen to"`
}

var (
	args = &cliArgs{
		Port: 6543,
	}
)

func getAPIServer() *restapi.Server {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	// TODO: remove logs
	// TODO: extract uuid logic and similar to the helper reducing dependencies
	// TODO: try to match origigal json as close as possible

	api := operations.NewForm3paymentsAPI(swaggerSpec)
	api.GetHealthHandler = operations.GetHealthHandlerFunc(GetHealth)
	api.GetPaymentsHandler = operations.GetPaymentsHandlerFunc(GetPayments)
	api.PostPaymentsHandler = operations.PostPaymentsHandlerFunc(CreatePayment)
	api.PatchPaymentsHandler = operations.PatchPaymentsHandlerFunc(UpdatePayment)
	api.GetPaymentsIDHandler = operations.GetPaymentsIDHandlerFunc(FetchPayment)
	api.DeletePaymentsIDHandler = operations.DeletePaymentsIDHandlerFunc(DeletePayment)
	server := restapi.NewServer(api)

	server.ConfigureAPI()

	return server
}

func main() {
	arg.MustParse(args)

	server := getAPIServer()
	defer server.Shutdown()

	server.Port = args.Port

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
