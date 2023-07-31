package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func VerifyAndChangeImage(errorRate float64, err error, latestTag string, previousTag string) {
	if errorRate < GlobalConfig.Kibana.ErrorRate {
		return
	}

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	deployment, err := clientset.AppsV1().Deployments(GlobalConfig.K8s.Namespace).Get(context.TODO(), GlobalConfig.K8s.DeploymentName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	imageTag := extractTagVersion(deployment.Spec.Template.Spec.Containers[0].Image)

	if imageTag == latestTag {
		newImage := GlobalConfig.K8s.RegistryImage + ":" + previousTag
		deployment.Spec.Template.Spec.Containers[0].Image = newImage

		_, updateErr := clientset.AppsV1().Deployments(GlobalConfig.K8s.Namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
		if updateErr != nil {
			panic(fmt.Errorf("failed to update Deployment: %w", updateErr))
		}

		fmt.Printf("Deployment '%s' in namespace '%s' updated. Rolling restart will be triggered.\n", GlobalConfig.K8s.DeploymentName, GlobalConfig.K8s.Namespace)
	}
}

func extractTagVersion(imageName string) string {
	imageParts := strings.Split(imageName, ":")
	if len(imageParts) > 1 {
		return imageParts[1]
	}

	return ""
}
