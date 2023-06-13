package archive

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ARTM2000/archive1/internal/config"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var parsedConfig Config

func archiveConfigPreProcess(configPath string) {
	finalConfigFile := strings.TrimSpace(configPath)
	if finalConfigFile == "" {
		// in this case, we set default configuration for config file
		home, err := homedir.Dir()
		if err != nil {
			log.Fatalln(err)
		}
		log.Default().Printf("no config file path received. looking at '%s' for '.archive1.yaml'", home)
		finalConfigFile = filepath.Join(home, ".archive1.yaml")
	}
	// check that finalConfigFile exists or not
	if _, err := os.Stat(finalConfigFile); os.IsNotExist(err) {
		log.Fatalf("no config file at '%s' found. error: %s\n", finalConfigFile, err.Error())
	}

	// parsing config file
	if err := config.Parse[Config](finalConfigFile, &parsedConfig); err != nil {
		log.Fatalf("error on reading configuration: %s", err.Error())
	}

	log.Default().Println("archive1 configuration:", parsedConfig.String())
	// validate received config
	if err := parsedConfig.Validate(); err != nil {
		log.Fatalf(err.Error())
	}

	log.Default().Println("configuration is valid")
}

var archiveCmd = &cobra.Command{
	Use:   "archive1",
	Short: "Archive1 server to store all agents files",
	Run: func(cmd *cobra.Command, _ []string) {
		configPath, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatalln(err.Error())
		}
		archiveConfigPreProcess(configPath)

		// run server with received configuration
		runServer(&parsedConfig)
	},
}

func init() {
	archiveCmd.Flags().StringP(
		"config",
		"c",
		"",
		"archive1 server configuration (default is $HOME/.archive1.yaml)",
	)
}

func CmdExecute() {
	if err := archiveCmd.Execute(); err != nil {
		log.Fatalln(err.Error())
	}
}
