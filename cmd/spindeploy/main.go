package main

import (
	"fmt"
	"os"
	"strings"

	spindeploy "github.com/robzienert/spin-deploy"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "spindeploy [target to deploy]",
	Short: "",
	Long:  ``,
	Run:   spindeploy.StartDeployHandler,
}

func init() {
	viper.SetConfigName(".spin")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("spindeploy")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func main() {
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("could not read config: " + err.Error())
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
