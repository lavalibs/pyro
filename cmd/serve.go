package cmd

import (
	"github.com/lavalibs/pyro/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		server := server.Server{}
		server.Serve(viper.Get("address").(string))
	},
}
