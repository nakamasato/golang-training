package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Skaffold struct {
	cmd *exec.Cmd
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fmt.Println("started")
	skaffold := &Skaffold{}
	skaffold.run(ctx, "--tail")
	defer skaffold.cleanup()
	k8sclient := getClientset()
	waitUntilPodIsReady(ctx, k8sclient, "test")
	listPod(ctx, k8sclient)
	fmt.Println("done")
}

func listPod(ctx context.Context, clientset *kubernetes.Clientset) {
	pods, err := clientset.CoreV1().Pods("default").List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatal("failed to get pods")
	}
	for i, pod := range pods.Items {
		fmt.Printf("Pod: [%d] %s\n", i, pod.GetName())
	}
}

func waitUntilPodIsReady(ctx context.Context, clientset *kubernetes.Clientset, name string) {
	for {
		pod, err := clientset.CoreV1().Pods("default").Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			time.Sleep(time.Second)
		}
		if pod.Status.Phase == "Running" {
			fmt.Println("pod is running")
			return
		} else {
			fmt.Println("waiting for pod to be ready")
			time.Sleep(time.Second)
		}
	}
}

func getClientset() *kubernetes.Clientset {
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
	return clientset
}

func (s *Skaffold) run(ctx context.Context, args ...string) {
	fmt.Println("run")
	args = append([]string{"run"}, args...)
	s.cmd = exec.CommandContext(
		ctx,
		"skaffold",
		args...,
	)
	s.cmd.Stdout = os.Stdout
	s.cmd.Stderr = os.Stderr
	s.cmd.Cancel = func() error {
		fmt.Println("cmd.Cancel is called") // not working
		return nil
	}
	s.cmd.Start() // Run in background
}

func (s *Skaffold) cleanup() error {
	fmt.Println("skaffold cleanup")
	fmt.Println("skaffold kill process")
	errKill := s.cmd.Process.Kill()
	s.cmd = exec.Command(
		"skaffold",
		"delete",
	)
	fmt.Println("skaffold delete")
	errRun := s.cmd.Run()
	if errKill != nil {
		return errKill
	} else if errRun != nil {
		return errRun
	}
	return nil
}
