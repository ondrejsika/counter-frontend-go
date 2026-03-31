package server

import (
	"github.com/ondrejsika/counter-frontend-go/internal/server"
)

type ServerOptions struct {
	VersionOverride string
}

func Server(options ...ServerOptions) {
	versionOverride := ""
	if len(options) > 0 {
		versionOverride = options[0].VersionOverride
	}
	server.Server(versionOverride)
}
