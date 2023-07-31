package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Kibana KibanaConfig `json:"kibana"`
	K8s    K8sConfig    `json:"k8s"`
	Docker DockerConfig `json:"docker"`
}

type KibanaConfig struct {
	Endpoint     string  `json:"endpoint"`
	ElasticIndex string  `json:"elasticIndex"`
	ErrorRate    float64 `json:"errorRate"`
}

type K8sConfig struct {
	Namespace      string `json:"namespace"`
	DeploymentName string `json:"deploymentName"`
	RegistryImage  string `json:"registryImage"`
}

type DockerConfig struct {
	RegistryURL string `json:"registryURL"`
	RepoName    string `json:"repoName"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

var GlobalConfig Config

func LoadConfig() error {
	file, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&GlobalConfig)
	if err != nil {
		return err
	}

	return nil
}