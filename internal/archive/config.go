package archive

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ARTM2000/archive1/internal/validate"
)

type Database struct {
	Host     string `mapstructure:"host" json:"host" validate:"required,hostname|ip"`
	Port     int    `mapstructure:"port" json:"port" validate:"required,number"`
	Username string `mapstructure:"username" json:"username" validate:"required"`
	Password string `mapstructure:"password" json:"password" validate:"required"`
	Name     string `mapstructure:"dbname" json:"dbname" validate:"required"`
	Zone     string `mapstructure:"timezone" json:"timezone" validate:"required"`
	SSLMode  bool   `mapstructure:"ssl_mode" json:"ssl_mode" validate:"omitempty,boolean"`
}

type Auth struct {
	JWTSecret     string        `mapstructure:"jwt_secret" json:"jwt_secret" validate:"required,min=10"`
	JWTExpireTime time.Duration `mapstructure:"jwt_expire_time" json:"jwt_expire_time" validate:"required"`
}

type Config struct {
	ServerPort *int     `mapstructure:"server_port" json:"server_port" validate:"omitempty,number"`
	ServerHost *string  `mapstructure:"server_host" json:"server_host" validate:"omitempty,hostname|ip"`
	Database   Database `mapstructure:"database" json:"database" validate:"required,dive"`
	Auth       Auth     `mapstructure:"auth" json:"auth" validate:"required,dive"`
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
