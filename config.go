package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

type Project struct {
	AccessToken string `yaml:"accessToken"`
	BaseUrl string `yaml:"baseUrl"`
	Repos []string
	MainBranch string `yaml:"mainBranch"`
}

type Config struct {
	AccessToken string `yaml:"accessToken"`
	BaseUrl string `yaml:"baseUrl"`
	Projects map[string]*Project
	MainBranch string `yaml:"mainBranch"`
}

func ResolveConfig(args *Args) (*Config, error) {
	filename, err := filepath.Abs(args.ConfigFile)
	
	if err != nil {
		return nil, err;
	}

	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err;
	}

	config := new(Config)
	err = yaml.Unmarshal(yamlFile, config)

	if err != nil {
		return nil, err;
	}

	if config.MainBranch == "" {
		config.MainBranch = "master"
	}

	for _, project := range config.Projects {
		if project.AccessToken == "" {
			project.AccessToken = config.AccessToken
		}

		if project.BaseUrl == "" {
			project.BaseUrl = config.BaseUrl
		}

		if project.MainBranch == "" {
			project.MainBranch = config.MainBranch
		}
	}

    return config, nil
} 
