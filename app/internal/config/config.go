package config

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

var DefaultSubstitutionPrefix = "~"

type Action struct {
	Alertname string
	Command   string
}

type Config struct {
	Version            string
	Actions            []Action
	SubstitutionPrefix string
}

type IConfigReader interface {
	Open(string) (io.ReadCloser, error)
	ReadAll(io.Reader) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

type ConfigReader struct{}

func (cr ConfigReader) Open(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

func (cr ConfigReader) ReadAll(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}

func (cr ConfigReader) Unmarshal(bytes []byte, v interface{}) error {
	return yaml.Unmarshal(bytes, v)
}

func Read(icr IConfigReader, path string) (Config, error) {
	var config Config

	file, err := icr.Open(path)
	if err != nil {
		slog.Error(err.Error())
		return config, err
	}
	defer file.Close()

	bytes, err := icr.ReadAll(file)
	if err != nil {
		slog.Error(err.Error())
		return config, err
	}

	if err := icr.Unmarshal(bytes, &config); err != nil {
		slog.Error(err.Error())
		return config, err
	}

	if config.SubstitutionPrefix == "" {
		slog.Info("no substitution prefix defined, using default: " + DefaultSubstitutionPrefix)
		config.SubstitutionPrefix = DefaultSubstitutionPrefix
	}
	return config, nil
}

func IsValid(config Config) (result bool) {
	if config.Version != "v1" {
		slog.Error("wrong config version: " + config.Version)
		return false
	}

	if len(config.Actions) == 0 {
		slog.Error("no actions defined")
		return false
	}

	for _, action := range config.Actions {
		if action.Alertname == "" {
			slog.Error("empty alertname in action: " + fmt.Sprint(action))
			return false
		}
		if action.Command == "" {
			slog.Error("empty command in action: " + fmt.Sprint(action))
			return false
		}
	}

	for i, action := range config.Actions {
		for j, action2 := range config.Actions {
			if i != j && action.Alertname == action2.Alertname {
				slog.Error("duplicate alertname in actions: " + fmt.Sprint(action) + " and " + fmt.Sprint(action2))
				return false
			}
		}
	}

	slog.Debug("config: " + fmt.Sprintf("%+v", config))

	return true
}
