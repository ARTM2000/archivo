package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ARTM2000/archive1/internal/validate"
	"github.com/robfig/cron/v3"
)

type File struct {
	Path     string `mapstructure:"path" json:"path" validate:"required,filepath"`
	Interval string `mapstructure:"interval" json:"interval" validate:"required"`
	Rotate   int64  `mapstructure:"rotate" json:"rotate" validate:"omitempty,required,number"`
	Filename string `mapstructure:"filename" json:"filename" validate:"omitempty,required,filename,alphanum"`
}

func (f *File) String() string {
	fileByte, _ := json.Marshal(f)
	return string(fileByte)
}

func (f *File) Validate() error {
	// block relative paths
	if !filepath.IsAbs(f.Path) {
		return fmt.Errorf("every paths should be absolute. invalid path: %s", f.Path)
	}

	// check file path existence
	if _, err := os.Stat(f.Path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not exists: %s", f.Path)
		}

		if !os.IsPermission(err) {
			return fmt.Errorf("file permission not granted: %s", f.Path)
		}
	}

	// check that received crontab is usable or not
	if _, err := cron.ParseStandard(f.Interval); err != nil {
		return fmt.Errorf("interval is invalid format: %s", err.Error())
	}

	if f.Rotate < 1 {
		return fmt.Errorf("rotate should be bigger than 0")
	}

	return nil
}

type Config struct {
	ArchiveServer string `mapstructure:"archive1_server" json:"archive1_server" validate:"required,url"`
	AgentName     string `mapstructure:"agent_name" json:"agent_name" validate:"required"`
	AgentKey      string `mapstructure:"agent_key" json:"-" validate:"required"`
	Files         []File `mapstructure:"files" json:"files" validate:"required,min=1,dive"`
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

	filenames := []string{}
	for _, file := range c.Files {
		if err := file.Validate(); err != nil {
			return err
		}
		if file.Filename != "" {
			filenames = append(filenames, file.Filename)
		}
	}

	isUnique, dup := validate.ValidateSliceParamUniqueness[string](filenames)
	if !isUnique {
		return fmt.Errorf("filename uniqueness violation. '%s' is a duplicate filename", *dup)
	}

	return nil
}
