package swagger

import (
	"github.com/swaggest/swgui/v5/static"
	"github.com/vearutop/statigz"
)

var staticServer = statigz.FileServer(static.FS)

const (
	__assets  = "{{ .BasePath }}"
	__favicon = "{{ .BasePath }}"
)
