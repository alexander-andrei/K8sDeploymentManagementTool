package gitops

import (
	"fmt"
	"io"
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

func ChangeAndCommit(applicationName string, previousTag string, imageName string) {
	// Replace these with your actual Git repository information
	gitCredentials := fmt.Sprintf("%s:%s@", config.GlobalConfig.GitOps.Username, config.GlobalConfig.GitOps.Token)

	repoURL := "https://" + gitCredentials + config.GlobalConfig.GitOps.Repo

	// Path to the file you want to change
	filePath := fmt.Sprintf("%s/%s-server-deployment.yaml", applicationName, applicationName)

	// Clone the repository to a temporary directory
	tmpDir := config.GlobalConfig.GitOps.TmpDir

	err := os.RemoveAll(tmpDir)
	if err != nil {
		fmt.Println("Error deleting contents:", err)
		return
	}

	// Directory is empty, perform a git clone
	cmd := exec.Command("git", "clone", repoURL, tmpDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error cloning repository:", err)
		return
	}

	// Write the new content to the file
	filePathInTmp := tmpDir + "/" + filePath

	yamlContent, err := os.ReadFile(filePathInTmp)
	if err != nil {
		fmt.Println("Error reading YAML file:", err)
		return
	}

	// Unmarshal the YAML content into a DeploymentSpec struct
	var deployment DeploymentSpec
	err = yaml.Unmarshal(yamlContent, &deployment)
	if err != nil {
		fmt.Println("Error unmarshalling YAML:", err)
		return
	}

	// Modify the image field
	deployment.Spec.Template.Spec.Containers[0].Image = utils.ReplaceImageTagVersion(imageName, previousTag)

	// Marshal the modified struct back to YAML
	updatedYAML, err := yaml.Marshal(deployment)
	err = os.WriteFile(filePathInTmp, []byte(updatedYAML), 0644)
	if err != nil {
		fmt.Println("Error writing new content:", err)
		return
	}

	// Commit and push the changes
	cmd = exec.Command("git", "add", filePath)
	cmd.Dir = tmpDir
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error adding file to git:", err)
		return
	}

	cmd = exec.Command("git", "commit", "-m", "Update file")
	cmd.Dir = tmpDir
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error committing changes:", err)
		return
	}

	cmd = exec.Command("git", "push")
	cmd.Dir = tmpDir
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error pushing changes:", err)
		return
	}

	fmt.Println("Changes pushed successfully!")
}

// isDirEmpty checks if a directory is empty.
func isDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
