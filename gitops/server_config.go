package gitops

import (
	"fmt"
	"io"
	"k8s/tool/config"
	"net/http"

	"gopkg.in/yaml.v2"
)

type Deployment struct {
	Kind     string `yaml:"kind"`
	Metadata struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Replicas int `yaml:"replicas"`
		Template struct {
			Spec struct {
				Containers []struct {
					Name  string `yaml:"name"`
					Image string `yaml:"image"`
				} `yaml:"containers"`
			} `yaml:"spec"`
		} `yaml:"template"`
	} `yaml:"spec"`
}

func GetServerConfigurationsFromGit(applicationName string) string {
	resp, err := http.Get(fmt.Sprintf("%s%s/%s-server-deployment.yaml", config.GlobalConfig.GitOps.RawRepo, applicationName, applicationName))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))

	var deployment Deployment
	err = yaml.Unmarshal(data, &deployment)
	if err != nil {
		panic(err)
	}

	return deployment.Spec.Template.Spec.Containers[0].Image
}
