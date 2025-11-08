// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	auth "Adornme/Auth"
	"Adornme/handlers"
	"Adornme/models"
	"Adornme/restapi/operations"
	"Adornme/restapi/operations/admin_products"
	"Adornme/restapi/operations/admin_users"
	"Adornme/restapi/operations/cart"
	"Adornme/restapi/operations/orders"
	"Adornme/restapi/operations/payments"
	"Adornme/restapi/operations/shipping"
	"Adornme/restapi/operations/users"
)

//go:generate swagger generate server --target ../../Adornme --name AdronmeCode --spec ../swagger/swagger.yaml --principal models.Principal

func configureFlags(api *operations.AdronmeCodeAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.AdronmeCodeAPI) http.Handler {
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
	api.BearerAuthAuth = func(token string) (*models.Principal, error) {
		// Validate token
		userID, err := auth.ValidateAccessToken(token)
		if err != nil {
			return nil, fmt.Errorf("invalid token: %w", err)
		}

		// Return a Principal object representing the logged-in user
		return &models.Principal{
			UserID: userID,
		}, nil
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()

	if api.CartAddItemToCartHandler == nil {
		api.CartAddItemToCartHandler = cart.AddItemToCartHandlerFunc(func(params cart.AddItemToCartParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation cart.AddItemToCart has not yet been implemented")
		})
	}
	if api.ShippingAddShippingAddressHandler == nil {
		api.ShippingAddShippingAddressHandler = shipping.AddShippingAddressHandlerFunc(func(params shipping.AddShippingAddressParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation shipping.AddShippingAddress has not yet been implemented")
		})
	}
	if api.CartClearCartHandler == nil {
		api.CartClearCartHandler = cart.ClearCartHandlerFunc(func(params cart.ClearCartParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation cart.ClearCart has not yet been implemented")
		})
	}
	if api.PaymentsConfirmPaymentHandler == nil {
		api.PaymentsConfirmPaymentHandler = payments.ConfirmPaymentHandlerFunc(func(params payments.ConfirmPaymentParams) middleware.Responder {
			return middleware.NotImplemented("operation payments.ConfirmPayment has not yet been implemented")
		})
	}
	if api.AdminProductsCreateProductHandler == nil {
		api.AdminProductsCreateProductHandler = admin_products.CreateProductHandlerFunc(func(params admin_products.CreateProductParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation admin_products.CreateProduct has not yet been implemented")
		})
	}
	if api.AdminProductsDeleteProductHandler == nil {
		api.AdminProductsDeleteProductHandler = admin_products.DeleteProductHandlerFunc(func(params admin_products.DeleteProductParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation admin_products.DeleteProduct has not yet been implemented")
		})
	}
	if api.ShippingDeleteShippingAddressHandler == nil {
		api.ShippingDeleteShippingAddressHandler = shipping.DeleteShippingAddressHandlerFunc(func(params shipping.DeleteShippingAddressParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation shipping.DeleteShippingAddress has not yet been implemented")
		})
	}
	if api.AdminUsersDeleteUserHandler == nil {
		api.AdminUsersDeleteUserHandler = admin_users.DeleteUserHandlerFunc(func(params admin_users.DeleteUserParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation admin_users.DeleteUser has not yet been implemented")
		})
	}
	if api.CartGetCartHandler == nil {
		api.CartGetCartHandler = cart.GetCartHandlerFunc(func(params cart.GetCartParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation cart.GetCart has not yet been implemented")
		})
	}
	if api.OrdersGetOrderHandler == nil {
		api.OrdersGetOrderHandler = orders.GetOrderHandlerFunc(func(params orders.GetOrderParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation orders.GetOrder has not yet been implemented")
		})
	}
	if api.PaymentsGetPaymentHandler == nil {
		api.PaymentsGetPaymentHandler = payments.GetPaymentHandlerFunc(func(params payments.GetPaymentParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation payments.GetPayment has not yet been implemented")
		})
	}
	if api.AdminUsersGetUserHandler == nil {
		api.AdminUsersGetUserHandler = admin_users.GetUserHandlerFunc(func(params admin_users.GetUserParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation admin_users.GetUser has not yet been implemented")
		})
	}

	if api.PaymentsInitiatePaymentHandler == nil {
		api.PaymentsInitiatePaymentHandler = payments.InitiatePaymentHandlerFunc(func(params payments.InitiatePaymentParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation payments.InitiatePayment has not yet been implemented")
		})
	}
	if api.OrdersListOrdersHandler == nil {
		api.OrdersListOrdersHandler = orders.ListOrdersHandlerFunc(func(params orders.ListOrdersParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation orders.ListOrders has not yet been implemented")
		})
	}
	if api.ShippingListShippingAddressesHandler == nil {
		api.ShippingListShippingAddressesHandler = shipping.ListShippingAddressesHandlerFunc(func(params shipping.ListShippingAddressesParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation shipping.ListShippingAddresses has not yet been implemented")
		})
	}
	if api.ShippingListShippingOptionsHandler == nil {
		api.ShippingListShippingOptionsHandler = shipping.ListShippingOptionsHandlerFunc(func(params shipping.ListShippingOptionsParams) middleware.Responder {
			return middleware.NotImplemented("operation shipping.ListShippingOptions has not yet been implemented")
		})
	}
	if api.AdminUsersListUsersHandler == nil {
		api.AdminUsersListUsersHandler = admin_users.ListUsersHandlerFunc(func(params admin_users.ListUsersParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation admin_users.ListUsers has not yet been implemented")
		})
	}
	if api.UsersLoginUserHandler == nil {
		api.UsersLoginUserHandler = users.LoginUserHandlerFunc(func(params users.LoginUserParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation users.LoginUser has not yet been implemented")
		})
	}
	if api.OrdersPlaceOrderHandler == nil {
		api.OrdersPlaceOrderHandler = orders.PlaceOrderHandlerFunc(func(params orders.PlaceOrderParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation orders.PlaceOrder has not yet been implemented")
		})
	}
	if api.PaymentsRefundPaymentHandler == nil {
		api.PaymentsRefundPaymentHandler = payments.RefundPaymentHandlerFunc(func(params payments.RefundPaymentParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation payments.RefundPayment has not yet been implemented")
		})
	}

	api.UsersRegisterUserHandler = users.RegisterUserHandlerFunc(handlers.RegisterUser)

	api.UsersGetUserProfileHandler = users.GetUserProfileHandlerFunc(handlers.GetUserProfile)

	api.UsersLoginUserHandler = users.LoginUserHandlerFunc(handlers.LoginUser)

	if api.UsersRequestPasswordResetHandler == nil {
		api.UsersRequestPasswordResetHandler = users.RequestPasswordResetHandlerFunc(func(params users.RequestPasswordResetParams) middleware.Responder {
			return middleware.NotImplemented("operation users.RequestPasswordReset has not yet been implemented")
		})
	}
	if api.UsersResetPasswordHandler == nil {
		api.UsersResetPasswordHandler = users.ResetPasswordHandlerFunc(func(params users.ResetPasswordParams) middleware.Responder {
			return middleware.NotImplemented("operation users.ResetPassword has not yet been implemented")
		})
	}
	if api.ShippingTrackShipmentHandler == nil {
		api.ShippingTrackShipmentHandler = shipping.TrackShipmentHandlerFunc(func(params shipping.TrackShipmentParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation shipping.TrackShipment has not yet been implemented")
		})
	}
	if api.CartUpdateCartItemHandler == nil {
		api.CartUpdateCartItemHandler = cart.UpdateCartItemHandlerFunc(func(params cart.UpdateCartItemParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation cart.UpdateCartItem has not yet been implemented")
		})
	}
	if api.AdminProductsUpdateProductHandler == nil {
		api.AdminProductsUpdateProductHandler = admin_products.UpdateProductHandlerFunc(func(params admin_products.UpdateProductParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation admin_products.UpdateProduct has not yet been implemented")
		})
	}
	if api.ShippingUpdateShippingAddressHandler == nil {
		api.ShippingUpdateShippingAddressHandler = shipping.UpdateShippingAddressHandlerFunc(func(params shipping.UpdateShippingAddressParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation shipping.UpdateShippingAddress has not yet been implemented")
		})
	}
	if api.AdminUsersUpdateUserHandler == nil {
		api.AdminUsersUpdateUserHandler = admin_users.UpdateUserHandlerFunc(func(params admin_users.UpdateUserParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation admin_users.UpdateUser has not yet been implemented")
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
