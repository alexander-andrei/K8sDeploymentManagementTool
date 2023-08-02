package gitops

import (
	"fmt"
	"k8s/tool/config"
	"k8s/tool/utils"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

type DeploymentSpec struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Replicas int `yaml:"replicas"`
		Selector struct {
			MatchLabels struct {
				App string `yaml:"app"`
			} `yaml:"matchLabels"`
		} `yaml:"selector"`
		Template struct {
			Metadata struct {
				Labels struct {
					App string `yaml:"app"`
				} `yaml:"labels"`
			} `yaml:"metadata"`
			Spec struct {
				Containers []struct {
					Name  string `yaml:"name"`
					Image string `yaml:"image"`
					Ports []struct {
						ContainerPort int `yaml:"containerPort"`
					} `yaml:"ports"`
				} `yaml:"containers"`
			} `yaml:"spec"`
		} `yaml:"template"`
	} `yaml:"spec"`
}

func ChangeAndCommit(applicationName string, previousTag string, imageName string) error {
	gitCredentials := fmt.Sprintf("%s:%s@", config.GlobalConfig.GitOps.Username, config.GlobalConfig.GitOps.Token)
	repoURL := "https://" + gitCredentials + config.GlobalConfig.GitOps.Repo
	filePath := fmt.Sprintf("%s/%s-server-deployment.yaml", applicationName, applicationName)
	tmpDir := config.GlobalConfig.GitOps.TmpDir

	err := os.RemoveAll(tmpDir)
	if err != nil {
		fmt.Println("Error deleting contents:", err)
		return err
	}

	cmd := exec.Command("git", "clone", repoURL, tmpDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error cloning repository:", err)
		return err
	}

	filePathInTmp := tmpDir + "/" + filePath

	yamlContent, err := os.ReadFile(filePathInTmp)
	if err != nil {
		fmt.Println("Error reading YAML file:", err)
		return err
	}

	var deployment DeploymentSpec
	err = yaml.Unmarshal(yamlContent, &deployment)
	if err != nil {
		fmt.Println("Error unmarshalling YAML:", err)
		return err
	}

	deployment.Spec.Template.Spec.Containers[0].Image = utils.ReplaceImageTagVersion(imageName, previousTag)

	updatedYAML, err := yaml.Marshal(deployment)
	err = os.WriteFile(filePathInTmp, []byte(updatedYAML), 0644)
	if err != nil {
		fmt.Println("Error writing new content:", err)
		return err
	}

	cmd = exec.Command("git", "add", filePath)
	cmd.Dir = tmpDir
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error adding file to git:", err)
		return err
	}

	cmd = exec.Command("git", "commit", "-m", "Update file")
	cmd.Dir = tmpDir
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error committing changes:", err)
		return err
	}

	cmd = exec.Command("git", "push")
	cmd.Dir = tmpDir
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error pushing changes:", err)
		return err
	}

	fmt.Println("Changes pushed successfully!")
	return nil
}
