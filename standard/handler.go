package standard

import "net/http"

// Handler handles swagger UI request.
type Handler interface {
	http.Handler
	GetBasePath() string
}
