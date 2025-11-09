package main

import (
	"log"
	"net/http"

	"Adornme/restapi"
	"Adornme/restapi/operations"

	"github.com/go-openapi/loads"
	httpSwagger "github.com/swaggo/http-swagger" // 👈 add this
)

func main() {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewAdronmeCodeAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	server.ConfigureAPI()

	server.Host = "0.0.0.0"
	server.Port = 8080
	server.TLSCertificate = ""
	server.TLSCertificateKey = ""

	// 👇 Add this Swagger UI route
	http.Handle("/docs/", httpSwagger.WrapHandler)
	go func() {
		log.Println("Swagger UI available at: http://localhost:8080/docs/index.html")
		log.Println("Server started at :8080 (HTTP)")
	}()

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
