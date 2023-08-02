package checker

import (
	"context"
	"flag"
	"fmt"
	jsonConfig "k8s/tool/config"
	"k8s/tool/utils"
	"path/filepath"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func VerifyAndChangeImage(errorRate float64, err error, latestTag string, previousTag string, deploymentName string, registryImage string) {
	if errorRate < jsonConfig.GlobalConfig.Kibana.ErrorRate {
		return
	}

	imageTag, deployment, clientset := getK8sDeoloymentInfo(deploymentName)

	if imageTag == latestTag {
		newImage := registryImage + ":" + previousTag
		deployment.Spec.Template.Spec.Containers[0].Image = newImage

		_, updateErr := clientset.AppsV1().Deployments(jsonConfig.GlobalConfig.K8s.Namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
		if updateErr != nil {
			panic(fmt.Errorf("failed to update Deployment: %w", updateErr))
		}

		fmt.Printf("Deployment '%s' in namespace '%s' updated. Rolling restart will be triggered.\n", deploymentName, jsonConfig.GlobalConfig.K8s.Namespace)
	}
}

func getK8sDeoloymentInfo(deploymentName string) (string, *v1.Deployment, *kubernetes.Clientset) {
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

	deployment, err := clientset.AppsV1().Deployments(jsonConfig.GlobalConfig.K8s.Namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	imageTag := utils.ExtractTagVersion(deployment.Spec.Template.Spec.Containers[0].Image)

	return imageTag, deployment, clientset
}
