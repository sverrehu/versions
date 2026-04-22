package config

import (
	_ "embed"
	"os"

	"go.yaml.in/yaml/v4"
)

//go:embed config.default.yaml
var defaultConfig []byte

type WebServer struct {
	Port         int `yaml:"port,omitempty"`
	CacheMinutes int `yaml:"cacheMinutes,omitempty"`
	CacheSize    int `yaml:"cacheSize,omitempty"`
}

type State struct {
	Filename string `yaml:"filename,omitempty"`
	Cache    Cache  `yaml:"cache,omitempty"`
}

type Cache struct {
	Releases         CacheSettings `yaml:"releases,omitempty"`
	CommitTimestamps CacheSettings `yaml:"commitTimestamps,omitempty"`
}

type CacheSettings struct {
	CacheMinutes int `yaml:"cacheMinutes,omitempty"`
	CacheSize    int `yaml:"cacheSize,omitempty"`
}

type Datasource struct {
	MaxReleases int          `yaml:"maxReleases,omitempty"`
	Credentials *Credentials `yaml:"credentials,omitempty"`
}

type Credentials struct {
	UserName string `yaml:"userName,omitempty"`
	Password string `yaml:"password,omitempty"`
	Token    string `yaml:"token,omitempty"`
}

type Datasources struct {
	GitHubReleasesDatasource *Datasource `yaml:"github-releases,omitempty"`
	GitHubTagsDatasource     *Datasource `yaml:"github-tags,omitempty"`
	GitLabReleasesDatasource *Datasource `yaml:"gitlab-releases,omitempty"`
	MavenDatasource          *Datasource `yaml:"maven,omitempty"`
	DockerhubDatasource      *Datasource `yaml:"dockerhub,omitempty"`
}

type Config struct {
	WebServer   WebServer    `yaml:"webServer,omitempty"`
	State       State        `yaml:"state,omitempty"`
	Datasources *Datasources `yaml:"datasources,omitempty"`
}

var cfg Config
var cfgLoaded = false

func Cfg() *Config {
	if !cfgLoaded {
		panic("config file not loaded")
	}
	return &cfg
}

func LoadConfig(fileName string) error {
	if cfgLoaded {
		panic("config file already loaded")
	}
	cfg = Config{}
	err := loadBytesInto(&cfg, defaultConfig)
	if err != nil {
		return err
	}
	if fileName != "" {
		err = loadFileInto(&cfg, fileName)
		if err != nil {
			return err
		}
	}
	cfgLoaded = true
	return nil
}

func loadFileInto(config *Config, filename string) error {
	yamlContents, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return loadBytesInto(config, yamlContents)
}

func loadBytesInto(config *Config, bytes []byte) error {
	return yaml.Unmarshal(bytes, config)
}
