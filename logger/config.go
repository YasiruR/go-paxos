package logger

import (
	"context"
	"github.com/tryfix/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Conf struct {
	ColorsEnabled bool   `yaml:"colors_enabled" default:"true"`
	LogLevel      string `yaml:"log_level" default:"ERROR"`
	FilePath      bool   `yaml:"file_path" default:"true"`
}

var config *Conf

func InitConfigs(ctx context.Context) {
	config = &Conf{}
	file, err := ioutil.ReadFile(`configs.yaml`)
	if err != nil {
		log.Fatal(`reading config file failed`)
	}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		log.Fatal(`unmarshalling configs failed`)
	}

	log.InfoContext(ctx, `log configurations initialized`)
}
