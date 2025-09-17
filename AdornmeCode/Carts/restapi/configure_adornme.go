// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"Adornme/AdornmeCode/Carts/restapi/operations"
	"Adornme/AdornmeCode/Carts/restapi/operations/cart"
	"Adornme/AdornmeCode/Carts/restapi/operations/orders"
)

//go:generate swagger generate server --target ..\..\Carts --name Adornme --spec ..\..\..\api\components\Carts\Carts.yaml --principal interface{} --exclude-main

func configureFlags(api *operations.AdornmeAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.AdornmeAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "Authorization" header is set
	if api.BearerAuthAuth == nil {
		api.BearerAuthAuth = func(token string) (interface{}, error) {
			return nil, errors.NotImplemented("api key auth (bearerAuth) Authorization from header param [Authorization] has not yet been implemented")
		}
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()

	if api.CartDeleteCartHandler == nil {
		api.CartDeleteCartHandler = cart.DeleteCartHandlerFunc(func(params cart.DeleteCartParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation cart.DeleteCart has not yet been implemented")
		})
	}
	if api.CartGetCartHandler == nil {
		api.CartGetCartHandler = cart.GetCartHandlerFunc(func(params cart.GetCartParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation cart.GetCart has not yet been implemented")
		})
	}
	if api.OrdersGetOrdersHandler == nil {
		api.OrdersGetOrdersHandler = orders.GetOrdersHandlerFunc(func(params orders.GetOrdersParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation orders.GetOrders has not yet been implemented")
		})
	}
	if api.OrdersGetOrdersIDHandler == nil {
		api.OrdersGetOrdersIDHandler = orders.GetOrdersIDHandlerFunc(func(params orders.GetOrdersIDParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation orders.GetOrdersID has not yet been implemented")
		})
	}
	if api.CartPostCartHandler == nil {
		api.CartPostCartHandler = cart.PostCartHandlerFunc(func(params cart.PostCartParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation cart.PostCart has not yet been implemented")
		})
	}
	if api.OrdersPostOrdersHandler == nil {
		api.OrdersPostOrdersHandler = orders.PostOrdersHandlerFunc(func(params orders.PostOrdersParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation orders.PostOrders has not yet been implemented")
		})
	}
	if api.CartPutCartHandler == nil {
		api.CartPutCartHandler = cart.PutCartHandlerFunc(func(params cart.PutCartParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation cart.PutCart has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
