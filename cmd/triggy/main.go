package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/lilic/triggy/pkg/triggy"
)

var (
	masterURL  string
	kubeconfig string
	image      string
)

func run() int {
	// Parse the flags.
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&image, "image", "", "The Docker image to deploy.")
	flag.Parse()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		fmt.Println("Error building kubeconfig: %s", err.Error())
		return 0
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		fmt.Println("Error building kubernetes clientset: %s", err.Error())
		return 1
	}

	t, err := triggy.New(triggy.Config{kubeClient})
	if err != nil {
		fmt.Println("Error creating new Trigger: %s", err.Error())
		return 1
	}

	err = t.Run(image)
	if err != nil {
		fmt.Println("Error while scalling deployment: %s", err.Error())
		return 1
	}

	return 0
}

func main() {
	os.Exit(run())
}
