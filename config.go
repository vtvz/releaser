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

func ResolveConfig() (Config) {
	filename, _ := filepath.Abs("./config.yaml")
	yamlFile, _ := ioutil.ReadFile(filename)
	
	var config Config
	_ = yaml.Unmarshal(yamlFile, &config)

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

    return config
} 
