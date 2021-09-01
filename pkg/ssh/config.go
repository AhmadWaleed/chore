package ssh

import (
	"fmt"
	"strings"
)

type Config struct {
	User   string
	Host   string
	Port   string
	RSA_ID string
}

func (c Config) IsLocalhost() bool {
	return c.Host == "localhost" || c.Host == "127.0.0.1"
}

func (c Config) Target() string {
	target := ""
	if c.User != "" {
		target += fmt.Sprintf("%s@%s", c.User, c.Host)
	} else {
		target += c.Host
	}

	if c.Port != "" {
		target += fmt.Sprintf(" -p %s", c.Port)
	}

	if c.RSA_ID != "" {
		target += fmt.Sprintf(" -i %s", c.RSA_ID)
	}

	return target
}

func ParseConfig(host string) Config {
	parts := strings.Split(host, "@")
	switch len(parts) {
	case 1:
		return Config{Host: parts[0]}
	case 2:
		return Config{User: parts[0], Host: parts[1]}
	case 3:
		return Config{User: parts[0], Host: parts[1], Port: parts[2]}
	case 4:
		return Config{User: parts[0], Host: parts[1], Port: parts[2], RSA_ID: parts[3]}
	default:
		return Config{}
	}
}
