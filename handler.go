package swagger

import (
	"github.com/thoohv5/swagger-ui/internal/swagger"
	"github.com/thoohv5/swagger-ui/internal/swagger/config"
	"github.com/thoohv5/swagger-ui/standard"
)

// New creates HTTP handler for Swagger UI.
func New(
	title, swaggerJSONPath string,
	basePath string,
) standard.Handler {
	return newHandler(title, swaggerJSONPath, basePath)
}

// NewWithOption creates configurable handler constructor.
func NewWithOption(handlerOpts ...HandlerOption) standard.Handler {
	opts := config.NewConfig()

	for _, o := range handlerOpts {
		o(opts)
	}

	return newHandlerWithConfig(opts)
}

// newHandlerWithConfig creates HTTP handler for Swagger UI.
func newHandlerWithConfig(config *config.Config) standard.Handler {
	return swagger.NewHandler(config, __assets, __favicon, staticServer)
}

// NewHandler creates HTTP handler for Swagger UI.
func newHandler(
	title, swaggerJSONPath string,
	basePath string,
) standard.Handler {
	return newHandlerWithConfig(&config.Config{
		Title:       title,
		SwaggerJSON: swaggerJSONPath,
		BasePath:    basePath,
	})
}
