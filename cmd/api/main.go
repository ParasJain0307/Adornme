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
	// Load the embedded swagger spec
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln("Error loading swagger spec:", err)
	}

	// Create API from spec
	api := operations.NewAdronmeCodeAPI(swaggerSpec)

	// Create the server
	server := restapi.NewServer(api)
	defer server.Shutdown()

	// Configure API
	server.ConfigureAPI()

	server.Host = "0.0.0.0"
	server.Port = 8080
	server.TLSCertificate = ""
	server.TLSCertificateKey = ""

	// Create a custom HTTP multiplexer
	mux := http.NewServeMux()

	// Serve your go-swagger generated API handlers
	mux.Handle("/", server.GetHandler())

	// Serve Swagger UI
	// Using relative path "/swagger.json" so it works inside EC2/Docker too
	mux.Handle("/docs/", httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"), // 👈 Swagger UI will load your API definition from here
		httpSwagger.DeepLinking(true),
	))

	log.Println("--------------------------------------------------")
	log.Println("🚀 Server running at: http://0.0.0.0:8080")
	log.Println("📘 Swagger UI available at: http://<EC2-PUBLIC-IP>:8080/docs/index.html")
	log.Println("--------------------------------------------------")

	// Start the HTTP server
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
