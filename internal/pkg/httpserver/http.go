package httpserver

import (
	"io"
	"net/http"

	"github.com/apoldev/trackchecker/internal/app/config"
	"github.com/apoldev/trackchecker/internal/app/restapi/restapi"
	"github.com/apoldev/trackchecker/internal/app/restapi/restapi/operations"
	trackhttp "github.com/apoldev/trackchecker/internal/app/track/delivery/http"
	"github.com/apoldev/trackchecker/pkg/logger"
	"github.com/bytedance/sonic"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
)

func NewOpenAPIServer(
	logger logger.Logger,
	trackHandlers *trackhttp.TrackHandler,
	cfg config.HTTPServer,
) *restapi.Server {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		logger.Fatal(err)
	}

	api := operations.NewTrackCheckerAppAPI(swaggerSpec)
	server := restapi.NewServer(api)

	server.Port = cfg.Port
	server.Host = cfg.Host
	server.EnabledListeners = []string{"http"}

	handler := configureAPI(api, logger, trackHandlers)
	server.SetHandler(handler)

	return server
}

func configureAPI(
	api *operations.TrackCheckerAppAPI,
	logger logger.Logger,
	trackHandlers *trackhttp.TrackHandler,
) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError
	api.Logger = func(s string, i ...interface{}) {
		logger.Infof(s, i...)
	}
	api.UseSwaggerUI()

	// Use Sonic instead of Standard json library.
	json := sonic.ConfigFastest

	api.JSONConsumer = runtime.ConsumerFunc(func(reader io.Reader, data interface{}) error {
		dec := json.NewDecoder(reader)
		dec.UseNumber() // preserve number formats
		return dec.Decode(data)
	})

	api.JSONProducer = runtime.ProducerFunc(func(writer io.Writer, data interface{}) error {
		enc := json.NewEncoder(writer)
		enc.SetEscapeHTML(false)
		return enc.Encode(data)
	})

	api.PostTrackHandler = operations.PostTrackHandlerFunc(trackHandlers.PostTrackingResultHandler)
	api.GetResultsHandler = operations.GetResultsHandlerFunc(trackHandlers.GetTrackingResultHandler)

	api.PreServerShutdown = func() {}
	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything,
// this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
