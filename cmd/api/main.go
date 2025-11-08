package main

import (
	"log"

	"Adornme/restapi"
	"Adornme/restapi/operations"

	"github.com/go-openapi/loads"
)

func main() {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewAdronmeCodeAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	// Force HTTP (skip TLS entirely)
	server.Host = "0.0.0.0"
	server.Port = 8080
	server.TLSCertificate = ""
	server.TLSCertificateKey = ""

	server.ConfigureAPI()

	log.Println("Server started at :8080 (HTTP)")
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
