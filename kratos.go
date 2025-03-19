package swagger

import (
	"fmt"
	"github.com/thoohv5/swagger-ui/internal/swagger/config"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

type openJsonFileHandler struct {
	Content []byte
}

func (h *openJsonFileHandler) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	_, _ = writer.Write(h.Content)
}

func (h *openJsonFileHandler) loadOpenApiFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	return content, err
}

func (h *openJsonFileHandler) LoadFile(filePath string) error {
	content, err := h.loadOpenApiFile(filePath)
	if err != nil {
		return err
	}

	h.Content = content
	return nil
}

type httpServer interface {
	HandlePrefix(prefix string, h http.Handler)
	Handle(path string, h http.Handler)
	HandleFunc(path string, h http.HandlerFunc)
}

func RegisterSwaggerUIServer[T httpServer](srv T, title, swaggerJSONPath string, basePath string) {
	handler := newHandler(title, swaggerJSONPath, basePath)
	srv.HandlePrefix(handler.GetBasePath(), handler)
}

func RegisterSwaggerUIServerWithOption[T httpServer](srv T, handlerOpts ...HandlerOption) {
	opts := config.NewConfig()

	for _, o := range handlerOpts {
		o(opts)
	}

	if opts.LocalOpenApiFile != "" {
		registerOpenApiLocalFileRouter(srv, opts)
	} else if len(opts.OpenApiData) != 0 {
		registerOpenApiMemoryDataRouter(srv, opts)
	}

	swaggerHandler := newHandlerWithConfig(opts)

	srv.HandlePrefix(swaggerHandler.GetBasePath(), swaggerHandler)
}

// var _openJsonFileHandler = &openJsonFileHandler{}

func registerOpenApiLocalFileRouter[T httpServer](srv T, cfg *config.Config) {
	var _openJsonFileHandler = &openJsonFileHandler{}
	err := _openJsonFileHandler.LoadFile(cfg.LocalOpenApiFile)
	if err == nil {
		pattern := strings.TrimRight(cfg.BasePath, "/") + "/openapi" + path.Ext(cfg.LocalOpenApiFile)
		cfg.SwaggerJSON = pattern
		srv.Handle(pattern, _openJsonFileHandler)
	} else {
		fmt.Println("load openapi file failed: ", err)
	}
}

func registerOpenApiMemoryDataRouter[T httpServer](srv T, cfg *config.Config) {
	var _openJsonFileHandler = &openJsonFileHandler{}
	_openJsonFileHandler.Content = cfg.OpenApiData
	pattern := strings.TrimRight(cfg.BasePath, "/") + "/openapi." + cfg.OpenApiDataType
	cfg.SwaggerJSON = pattern
	srv.Handle(pattern, _openJsonFileHandler)
	cfg.OpenApiData = nil
}
