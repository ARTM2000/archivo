package agent

import (
	"encoding/json"
)

type File struct {
	Path     string `mapstructure:"path" json:"path"`
	Interval string `mapstructure:"interval" json:"interval"`
}

func (f *File) String() string {
	fileByte, _ := json.Marshal(f)
	return string(fileByte)
}

type Config struct {
	ArchiveServer string `mapstructure:"archive1_server" json:"archive1_server"`
	ArchiveKey    string `mapstructure:"archive1_key" json:"-"`
	Files         []File `mapstructure:"files" json:"files"`
}

func (c *Config) String() string {
	configBytes, _ := json.Marshal(c)
	return string(configBytes)
}
