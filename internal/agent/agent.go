package agent

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ARTM2000/archive1/internal/config"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var configFile *string
var parsedConfig Config

var agentCmd = &cobra.Command{
	Use:   "agent1",
	Short: "Archive1 Agent to send specified files to Archive1 server",
	Run: func(cmd *cobra.Command, args []string) {
		finalConfigFile := *configFile
		if strings.TrimSpace(finalConfigFile) == "" {
			// in this case, we set default configuration for config file
			log.Default().Println("no config file path received. looking at $HOME for '.agent1.yaml'")
			home, err := homedir.Dir()
			if err != nil {
				log.Fatalln(err)
			}
			finalConfigFile = filepath.Join(home, ".agent1.yaml")
		}
		// check that finalConfigFile exists or not
		if _, err := os.Stat(finalConfigFile); os.IsNotExist(err) {
			log.Fatalf("no config file at '%s' found\nerror: %s\n", finalConfigFile, err.Error())
		}

		config.ParseYaml[Config](finalConfigFile, &parsedConfig)
		log.Default().Println("agent1 configuration:\n", parsedConfig.String())

	},
}

func init() {
	configFile = agentCmd.Flags().StringP(
		"config",
		"c",
		"",
		"path of agent1 config yaml file (default to $HOME/.agent1.yaml)",
	)
}

func CmdExecute() {
	if err := agentCmd.Execute(); err != nil {
		log.Fatalln(err.Error())
	}
}
