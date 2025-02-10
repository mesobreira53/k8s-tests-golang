package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/clientcmd"

	// To Manage custom resources
	kclient "github.com/mesobreira53/k8s-tests-golang/k8sclient"
	vault "github.com/mesobreira53/k8s-tests-golang/vaultclient"

	"k8s.io/client-go/dynamic"
)

func main() {

	fmt.Println("Starting Kubernetes connection...")
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("error getting user home dir: %v\n", err)
		os.Exit(1)
	}
	defaultKubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")

	kubeConfigPath := flag.String("kconfig", defaultKubeConfigPath, "Kubeconfig path")
	flag.Parse()
	fmt.Println("Get Kubernetes pods")

	/*	fmt.Println("Using kubeconfig ", *kubeConfigPath)
		kubeConfig, err := clientcmd.BuildConfigFromFlags("", *kubeConfigPath)
		if err != nil {
			fmt.Printf("error getting Kubernetes config: %v\n", err)
			os.Exit(1)
		}

		clientset, err := kubernetes.NewForConfig(kubeConfig)
		if err != nil {
			fmt.Printf("error getting Kubernetes clientset: %v\n", err)
			os.Exit(1)
		}*/

	k8sClient, err := kclient.NewClient(*kubeConfigPath)
	if err != nil {
		os.Exit(1)
	}

	// Pods in Kube-system
	// pods, err := clientset.CoreV1().Pods("kube-system").List(context.Background(), v1.ListOptions{})

	// Pods in all namespaces
	pods, err := k8sClient.Clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("error getting pods: %v\n", err)
		os.Exit(1)
	}
	for _, pod := range pods.Items {
		fmt.Printf("Pod name: %s from %s\n", pod.Name, pod.Namespace)
	}

	// Define the Secret you want to generate
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-secret", // Name of the secret
			Namespace: "default",   // Namespace where the secret will be created
		},
		Type: v1.SecretTypeOpaque, // Type of secret (Opaque is the default type)
		Data: map[string][]byte{
			"username": []byte("admin"),    // Key "username" with value "admin"
			"password": []byte("P@ssw0rd"), // Key "password" with value "P@ssw0rd"
		},
	}

	//Create the secret in the cluster
	secretsClient := k8sClient.Clientset.CoreV1().Secrets("default") // Specify the namespace where the secret should go
	currentSecret, err := secretsClient.Get(context.Background(), "my-secret", metav1.GetOptions{})
	if err != nil {
		log.Fatalf("Failed to get the secret. Looks like it does not exist")
		createdSecret, err := secretsClient.Create(context.Background(), secret, metav1.CreateOptions{})
		if err != nil {
			log.Fatalf("Failed to create secret: %v", err)
		}
		fmt.Printf("Secret %s created successfully!\n", createdSecret.Name)
	} else {
		fmt.Printf("Secret %s already exist\n", currentSecret.Name)
		updatedSecret, err := secretsClient.Update(context.Background(), secret, metav1.UpdateOptions{})
		if err != nil {
			log.Fatalf("Failed to update secret: %v", err)
		}
		fmt.Printf("Secret %s updated successfully!\n", updatedSecret.Name)
		//return
	}

	// Custome Resource
	// Build Kubernetes dynamic client
	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfigPath)
	if err != nil {
		panic(err)
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	// Read my Custom Resource
	gvr := schema.GroupVersionResource{
		Group:    "k8s.mes.local",
		Version:  "v1",
		Resource: "webresources", // This should be the plural name of your CRD
	}

	crName := "my-web-server"
	namespace := "default"
	customResource, err := dynamicClient.Resource(gvr).Namespace(namespace).Get(context.TODO(), crName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	// Convert Unstructured CR to JSON and then to Struct
	crData, _ := json.Marshal(customResource.Object)
	var myCR WebResource
	err = json.Unmarshal(crData, &myCR)
	if err != nil {
		log.Fatal(err)
	}
	// Print the CR Spec values
	fmt.Println("Custom Resource Details:")
	fmt.Printf("Name: %s\n", myCR.ObjectMeta.Name)
	fmt.Printf("Replicas: %d\n", myCR.Spec.Replicas)
	fmt.Printf("Image: %s\n", myCR.Spec.Image)

	// Check Secrets
	toSync, err := k8sClient.K8SReadSecret("mes-local-cert-tls", "cert-manager")

	if err != nil {
		fmt.Println("Error checking Secrets")
		os.Exit(1)
	}
	if toSync {
		fmt.Println("Temos de sincronizar")
	}

	// Vault integration
	fmt.Println("Connecting to Vault...")
	v := vault.NewVault("http://192.168.1.200:8200", "hvs")
	v.VaultSyncSecret()

}
