package cmd

import (
	"github.com/lavalibs/pyro/http"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)


var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		server := http.Server{}
		server.Serve(viper.Get("address").(string))
	},
}
