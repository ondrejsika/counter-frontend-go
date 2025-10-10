package root

import (
	"github.com/ondrejsika/counter-frontend-go/pkg/server"
	"github.com/ondrejsika/counter-frontend-go/version"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "counter-frontend-go",
	Short: "counter-frontend-go, " + version.Version,
	Run: func(c *cobra.Command, args []string) {
		server.Server()
	},
}
