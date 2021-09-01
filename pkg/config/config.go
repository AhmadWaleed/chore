package config

import (
	"fmt"

	"github.com/ahmadwaleed/chore/pkg/ssh"
)

type Bucket struct {
	Name  string   `yaml:"name"`
	Tasks []string `yaml:"tasks"`
}

type EnvVar map[string]string

func (ev EnvVar) ToStringVars() []string {
	var env []string
	for k, v := range ev {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	return env
}

type Task struct {
	Name     string   `yaml:"name"`
	Env      EnvVar   `yaml:"env"`
	Commands []string `yaml:"run"`
	Hosts    []ssh.Config
}

type Config struct {
	Servers []string  `yaml:"servers"`
	Tasks   []*Task   `yaml:"tasks"`
	Buckets []*Bucket `yaml:"buckets"`
	Vars    EnvVar    `yaml:"vars"`
}

func NewConfig(path string) (Config, error) {
	conf, err := NewParser(path).Parse()
	if err != nil {
		return conf, err
	}

	var configs []ssh.Config
	for _, h := range conf.Servers {
		configs = append(configs, ssh.ParseConfig(h))
	}

	for _, t := range conf.Tasks {
		t.Hosts = configs
		if conf.Vars == nil {
			continue
		}

		if t.Env != nil {
			for k, v := range conf.Vars {
				t.Env[k] = v
			}
		} else {
			t.Env = conf.Vars
		}
	}

	return conf, nil
}
