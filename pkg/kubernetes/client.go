package kubernetes

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func GetClientSet(kubeconfig string) *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return cs
}
