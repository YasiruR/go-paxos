package domain

import (
	"context"
	"github.com/tryfix/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Conf struct {
	LeaderTimeout  int64 `yaml:"leader_http_timeout"`
	ReplicaTimeout int64 `yaml:"replica_http_timeout"`
}

var Config *Conf

func SetConfigs(ctx context.Context) {
	Config = &Conf{}
	file, err := ioutil.ReadFile(`configs.yaml`)
	if err != nil {
		log.Fatal(`reading config file failed`)
	}

	err = yaml.Unmarshal(file, Config)
	if err != nil {
		log.Fatal(`unmarshalling configs failed`)
	}

	log.InfoContext(ctx, `domain configurations initialized`)
}
