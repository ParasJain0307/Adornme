package main

import (
	"log"
	"net/http"

	"Adornme/restapi"
	"Adornme/restapi/operations"

	"github.com/go-openapi/loads"
	httpSwagger "github.com/swaggo/http-swagger"
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

	// 👇 Add swagger UI and swagger.json routing
	swaggerHandler := httpSwagger.WrapHandler
	handler := http.NewServeMux()
	handler.Handle("/docs/", swaggerHandler)
	handler.Handle("/swagger.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(restapi.FlatSwaggerJSON)
	}))
	handler.Handle("/", server.GetHandler())

	server.SetHandler(handler)

	server.Host = "0.0.0.0"
	server.Port = 8080
	server.TLSCertificate = ""
	server.TLSCertificateKey = ""

	log.Println("🚀 Server started on :8080")
	log.Println("📘 Swagger UI: http://<EC2-IP>:8080/docs/index.html")
	log.Println("📄 Swagger JSON: http://<EC2-IP>:8080/swagger.json")

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
