package spindeploy

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/robzienert/spin-deploy/config"
	"github.com/robzienert/spin-deploy/spinnaker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func StartDeployHandler(cmd *cobra.Command, args []string) {
	clientConfig, err := spinnaker.NewConfig()
	if err != nil {
		log.Println("failed creating spinnaker client config: " + err.Error())
		os.Exit(1)
	}

	client, err := spinnaker.New(*clientConfig)
	if err != nil {
		log.Println("failed creating spinnaker client: " + err.Error())
		os.Exit(1)
	}

	targetName := args[0]
	target, err := config.GetTarget(targetName)
	if err != nil {
		log.Println("failed getting template source: " + err.Error())
		os.Exit(1)
	}

	req := spinnaker.StartPipelineRequest{
		Application: getAppName(viper.GetString("metadata.app")),
		Name:        getPipelineName(targetName),
		Template: spinnaker.DCDPipelineTemplate{
			Source: target.Template,
		},
		Configuration: spinnaker.DCDPipelineConfig{},
	}

	ref, err := client.StartPipeline(req)
	if err != nil {
		log.Println("error starting pipeline: " + err.Error())
		os.Exit(1)
	}

	// TODO tailing, etc etc

	fmt.Println(ref)
	os.Exit(0)
}

func getPipelineName(target string) string {
	return fmt.Sprintf("DCD Deploy to %s", target)
}

func getAppName(appName string) string {
	if appName == "" {
		appName = filepath.Base(filepath.Dir(os.Args[0]))
	}
	return strings.Replace(appName, "-", "", -1)
}
