package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Prefix         string            `yaml:"prefix"`
	ExternalSecret ExternalSecretOpt `yaml:"externalSecret"`
	ExportOpts     []ExportOpt       `yaml:"exportOpts"`
	Tags           map[string]string `yaml:"tags"`
}

type ExternalSecretOpt struct{}

type ExportOpt struct {
	SecretName      string `yaml:"secretName"`
	SecretNamespace string `yaml:"secretNamespace"`
	SecretKey       string `yaml:"secretKey"`
	DecodeBase64    bool   `yaml:"decodeBase64"`
}

func Load(p string) (*Config, error) {
	ret := &Config{}
	buf, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(buf, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
