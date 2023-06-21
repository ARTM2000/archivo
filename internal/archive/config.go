package archive

import (
	"encoding/json"
	"fmt"
	"os"
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

type FileStore struct {
	Mode       string        `mapstructure:"mode" json:"mode" validate:"required,eq=disk|eq=minio"`
	DiskConfig DiskFileStore `mapstructure:"disk_config" json:"disk_config" validate:"omitempty,required,dive"`
}

func (fs *FileStore) Validate() error {
	switch fs.Mode {
	case "disk":
		if fs.DiskConfig == (DiskFileStore{}) {
			return fmt.Errorf("for disk mode, disk_config is required")
		}
		pathD, err := os.Stat(fs.DiskConfig.Path)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("path '%s' not exists", fs.DiskConfig.Path)
			}
			if os.IsPermission(err) {
				return fmt.Errorf("access denied to access path '%s'", fs.DiskConfig.Path)
			}
			if !pathD.IsDir() {
				return fmt.Errorf("path is not valid directory")
			}
		}
	case "minio":
		// todo: add support for minio
		return fmt.Errorf("minio is not supported yet")
	default:
		return fmt.Errorf("file store mode not defined. mode: '%s'", fs.Mode)
	}
	return nil
}

type DiskFileStore struct {
	Path string `mapstructure:"path" json:"path" validate:"required,dir"`
}

type Config struct {
	ServerPort *int      `mapstructure:"server_port" json:"server_port" validate:"omitempty,number"`
	ServerHost *string   `mapstructure:"server_host" json:"server_host" validate:"omitempty,hostname|ip"`
	Database   Database  `mapstructure:"database" json:"database" validate:"required,dive"`
	Auth       Auth      `mapstructure:"auth" json:"auth" validate:"required,dive"`
	FileStore  FileStore `mapstructure:"file_store" json:"file_store" validate:"required"`
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

	fileStoreErr := c.FileStore.Validate()
	if fileStoreErr != nil {
		return fmt.Errorf("file store config got error. %s", fileStoreErr.Error())
	}

	return nil
}
