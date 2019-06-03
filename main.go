package main

import (
	"log"

	"github.com/alexflint/go-arg"
	"github.com/go-openapi/loads"
	"github.com/pawmart/form3-payments/config"
	"github.com/pawmart/form3-payments/storage"

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

	cfg := new(config.Config).LoadConfiguration()
	s := &storage.Storage{Config: cfg.Db}
	h := &PaymentsHandler{storage: s}

	api := operations.NewForm3paymentsAPI(swaggerSpec)
	api.GetHealthHandler = operations.GetHealthHandlerFunc(h.GetHealth)
	api.GetPaymentsHandler = operations.GetPaymentsHandlerFunc(h.GetPayments)
	api.PostPaymentsHandler = operations.PostPaymentsHandlerFunc(h.CreatePayment)
	api.PatchPaymentsHandler = operations.PatchPaymentsHandlerFunc(h.UpdatePayment)
	api.GetPaymentsIDHandler = operations.GetPaymentsIDHandlerFunc(h.FetchPayment)
	api.DeletePaymentsIDHandler = operations.DeletePaymentsIDHandlerFunc(h.DeletePayment)

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
