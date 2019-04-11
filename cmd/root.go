package cmd


import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "pyro",
	Short: "Pyro is a distributable state container for Lavalink.",
}

var configFile string

func init() {
	rootCmd.AddCommand(serveCmd)

	cobra.OnInitialize(initConfig)
	serveCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Set the location of the config file")
	serveCmd.PersistentFlags().StringP("addr", "a", ":80", "Bind to the specified address")
	viper.BindPFlag("address", serveCmd.PersistentFlags().Lookup("addr"))
}

func initConfig() {
  // Don't forget to read config either from cfgFile or from home directory!
  if configFile != "" {
    // Use config file from the flag.
    viper.SetConfigFile(configFile)
  } else {
    // Find home directory.
    home, err := homedir.Dir()
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    // Search config in home directory with name ".pyro" (without extension).
    viper.AddConfigPath(home)
    viper.SetConfigName(".pyro")
  }

  if err := viper.ReadInConfig(); err != nil {
    fmt.Println("Can't read config:", err)
    os.Exit(1)
  }
}

// Execute this CLI app
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
