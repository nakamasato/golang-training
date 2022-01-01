package main

import (
	"context"
	"fmt"
	"os"
	"path"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/kubectl/pkg/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	corev1 "k8s.io/api/core/v1"
)

func main() {
	k8sClient, err := getK8sClient()
	if err != nil {
		fmt.Println("failed to get k8sClient")
		return
	}

	getPod(k8sClient, "test-pod", "default")
	getPod(k8sClient, "kube-controller-manager-test-control-plane", "kube-system")
}

func getPod(k8sClient client.Client, name, namespace string) {
	pod := &corev1.Pod{}
	err := k8sClient.Get(context.TODO(), client.ObjectKey{Namespace: namespace, Name: name}, pod)
	if err != nil {
		if errors.IsNotFound(err) {
			fmt.Printf("%s pod doesn't exist in %s namespace\n", name, namespace)
			return
		} else {
			fmt.Println(err)
			return
		}
	}
	fmt.Printf("%s pod exists in %s namespace\n", name, namespace)
}

func getK8sClient() (client.Client, error) {
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	os.Setenv("KUBECONFIG", path.Join(mydir, "kubeconfig"))
	cfg, err := config.GetConfigWithContext("kind-test")
	if err != nil {
		return nil, err
	}

	k8sClient, err := client.New(cfg, client.Options{Scheme: scheme.Scheme})
	return k8sClient, err
}
