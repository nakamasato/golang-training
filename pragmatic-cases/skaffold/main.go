package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Skaffold struct{}

func main() {
	fmt.Println("started")
	skaffold := &Skaffold{}
	skaffold.run()
	listPod()
	skaffold.delete()
	fmt.Println("done")
}

func listPod() {
	var kubeconfig string
	kubeconfig, ok := os.LookupEnv("KUBECONFIG")
	if !ok {
		kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal("failed to get pods")
	}
	for i, pod := range pods.Items {
		fmt.Printf("Pod: [%d] %s\n", i, pod.GetName())
	}
}

func (s *Skaffold) run() {
	s.execute("run")
}

func (s *Skaffold) delete() {
	s.execute("delete")
}

func (s *Skaffold) execute(args ...string) {
	cmd := exec.Command(
		"skaffold",
		args...,
	)
	// cmd.Dir = "."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal("failed to run skaffold dev.", err)
	}
	fmt.Println("skaffold run completed")
}
