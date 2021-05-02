package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	Routes[] struct {
		PathPrefix string `yaml:"path_prefix"`
		Backend string `yaml:"backend"`
	} `yaml:"routes"`
	DefaultResponse struct {
		Body string `yaml:"body"`
		StatusCode int `yaml:"status_code"`
	} `yaml:"default_response"`
	Backends[] Backend `yaml:"backends"`
}

type Backend struct {
	Name string `yaml:"name"`
	MatchLabels struct {
		AppName string `yaml:"app_name"`
		Env string `yaml:"env"`
	} `yaml:"match_labels"`
}

func ParseConfig(configFilePath string) (config *Config){
	InfoLogger.Println("Reading configuration.")
	config = &Config{}
	configFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		ErrorLogger.Printf("Unable to read config file: %s \n", configFilePath)
		os.Exit(1)
	}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		ErrorLogger.Printf("Unable to Parse Yaml File : %s. \n", configFile)
		os.Exit(1)
	}
	return config
}


