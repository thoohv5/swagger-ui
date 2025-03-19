package swagger

import (
	"encoding/json"
	"github.com/thoohv5/swagger-ui/internal/swagger/config"
	"github.com/thoohv5/swagger-ui/internal/swagger/tpl"
	"github.com/thoohv5/swagger-ui/standard"
	"html/template"
	"net/http"
	"strings"
)

// handler handles swagger UI request.
type handler struct {
	*config.Config

	ConfigJson template.JS

	tpl          *template.Template
	staticServer http.Handler
}

// NewHandler returns a HTTP handler for swagger UI.
func NewHandler(
	config *config.Config,
	assets, favicon string,
	staticServer http.Handler,
) standard.Handler {
	config.BasePath = strings.TrimSuffix(config.BasePath, "/") + "/"

	h := &handler{
		Config: config,
	}

	j, err := json.Marshal(h.Config)
	if err != nil {
		panic(err)
	}

	h.ConfigJson = template.JS(j) //nolint:gosec // Data is well-formed.

	h.tpl, err = template.New("index").Parse(tpl.Index(assets, favicon, config))
	if err != nil {
		panic(err)
	}

	if staticServer != nil {
		h.staticServer = http.StripPrefix(h.BasePath, staticServer)
	}

	return h
}

// ServeHTTP implements http.Handler interface to handle swagger UI request.
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.TrimSuffix(r.URL.Path, "/") != strings.TrimSuffix(h.BasePath, "/") && h.staticServer != nil {
		h.staticServer.ServeHTTP(w, r)

		return
	}

	w.Header().Set("Content-Type", "text/html")

	if err := h.tpl.Execute(w, h); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *handler) GetBasePath() string {
	return h.BasePath
}
