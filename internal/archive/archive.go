package archive

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ARTM2000/archivo/internal/config"
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
		log.Default().Printf("no config file path received. looking at '%s' for '.archivo.yaml'", home)
		finalConfigFile = filepath.Join(home, ".archivo.yaml")
	}
	// check that finalConfigFile exists or not
	if _, err := os.Stat(finalConfigFile); os.IsNotExist(err) {
		log.Fatalf("no config file at '%s' found. error: %s\n", finalConfigFile, err.Error())
	}

	// parsing config file
	if err := config.Parse[Config](finalConfigFile, &parsedConfig); err != nil {
		log.Fatalf("error on reading configuration: %s", err.Error())
	}

	log.Default().Println("archivo configuration:", parsedConfig.String())
	// validate received config
	if err := parsedConfig.Validate(); err != nil {
		log.Fatalf(err.Error())
	}

	log.Default().Println("configuration is valid")
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate Archivo configuration",
	Run: func(cmd *cobra.Command, _ []string) {
		configPath, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatalln(err.Error())
		}
		archiveConfigPreProcess(configPath)
	},
}

var archiveCmd = &cobra.Command{
	Use:   "archivo",
	Short: "Archivo server to store all agents files",
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
		"archivo server configuration (default is $HOME/.archivo.yaml)",
	)

	validateCmd.Flags().StringP(
		"config",
		"c",
		"",
		"archivo server configuration (default is $HOME/.archivo.yaml)",
	)
}

func CmdExecute() {
	archiveCmd.AddCommand(validateCmd)
	if err := archiveCmd.Execute(); err != nil {
		log.Fatalln(err.Error())
	}
}
