package k8sclient

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sConnection struct {
	Clientset *kubernetes.Clientset
}

type K8sClient interface {
	K8SReadSecret(secretName, namespace string) (bool, error)
}

func NewClient(kconfig string) (*K8sConnection, error) {
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kconfig)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		fmt.Printf("error getting Kubernetes clientset: %v\n", err)
		return nil, err
	}
	return &K8sConnection{Clientset: clientset}, nil
}

func (c *K8sConnection) K8SReadSecret(secretName string, nameSpace string) (bool, error) {

	secretsClient := c.Clientset.CoreV1().Secrets(nameSpace) // Specify the namespace where the secret should go
	currentSecret, err := secretsClient.Get(context.Background(), secretName, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("Cannot read the Secret %s in namespace %s. Looks like it does not exists.", secretName, nameSpace)
		return false, err
	}
	fmt.Println(currentSecret)
	secret_generation := currentSecret.Generation
	fmt.Printf("Current generation of %s in %s is %d ", secretName, nameSpace, secret_generation)
	return true, nil
}
