package agent

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ARTM2000/archivo/internal/config"
	"github.com/ARTM2000/archivo/internal/processmng"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// var configFile *string
var parsedConfig Config

func agentConfigPreProcess(configPath string) {
	finalConfigFile := strings.TrimSpace(configPath)
	if finalConfigFile == "" {
		// in this case, we set default configuration for config file
		home, err := homedir.Dir()
		if err != nil {
			log.Fatalln(err)
		}
		log.Default().Printf("no config file path received. looking at '%s' for '.agent.yaml'", home)
		finalConfigFile = filepath.Join(home, ".agent.yaml")
	}
	// check that finalConfigFile exists or not
	if _, err := os.Stat(finalConfigFile); os.IsNotExist(err) {
		log.Fatalf("no config file at '%s' found. error: %s\n", finalConfigFile, err.Error())
	}

	// parsing config file
	if err := config.Parse[Config](finalConfigFile, &parsedConfig); err != nil {
		log.Fatalf("error on reading configuration: %s", err.Error())
	}

	log.Default().Println("agent configuration:", parsedConfig.String())
	// validate received config
	if err := parsedConfig.Validate(); err != nil {
		log.Fatalf(err.Error())
	}

	log.Default().Println("configuration is valid")
}

var validateAgentCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate Archivo Agent configuration",
	Run: func(cmd *cobra.Command, _ []string) {
		configPath, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatalf(err.Error())
		}
		agentConfigPreProcess(configPath)
	},
}

var agentCmd = &cobra.Command{
	Use:   "agent1",
	Short: "Archivo Agent to send specified files to Archivo server",
	Run: func(cmd *cobra.Command, _ []string) {
		configPath, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatalf(err.Error())
		}
		agentConfigPreProcess(configPath)
		agCron, err := registerCronJobs(&parsedConfig)
		if err != nil {
			log.Fatalf(err.Error())
		}

		eCh := make(chan int)
		go processmng.OnInterrupt(func() {
			agCron.Stop()
			eCh <- 1
		})

		// in order to keep app running
		<-eCh
	},
}

func init() {
	agentCmd.Flags().StringP(
		"config",
		"c",
		"",
		"path of agent1 config yaml file (default to $HOME/.agent.yaml)",
	)

	validateAgentCmd.Flags().StringP(
		"config",
		"c",
		"",
		"path of agent1 config yaml file (default to $HOME/.agent.yaml)",
	)
}

func CmdExecute() {
	agentCmd.AddCommand(validateAgentCmd)
	if err := agentCmd.Execute(); err != nil {
		log.Fatalln(err.Error())
	}
}
