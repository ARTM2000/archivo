package archive

import (
	"encoding/json"
	"fmt"

	"github.com/ARTM2000/archive1/internal/validate"
)

type Config struct {
	ServerPort *int    `mapstructure:"server_port" json:"server_port" validate:"omitempty,number"`
	ServerHost *string `mapstructure:"server_host" json:"server_host" validate:"omitempty,hostname|ip"`
}

func (c *Config) String() string {
	configBytes, _ := json.Marshal(c)
	return string(configBytes)
}

func (c *Config) Validate() error {
	errors, ok := validate.ValidateStruct[Config](c)
	if !ok {
		return fmt.Errorf("configuration validation error: %s", errors[0].Message)
	}

	return nil
}
