// Package kubectl
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package kubectl

import (
	"context"
	apiv1 "k8s.io/api/core/v1"
	errosv1 "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func CreateNamespace(kubeconfig string, namespace string) error {

	cs := GetClientSet(kubeconfig)
	ns := &apiv1.Namespace{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Namespace",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
			Labels: map[string]string{
				"name": namespace,
			},
		},
	}

	_, err := cs.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{})
	if err != nil && !errosv1.IsAlreadyExists(err) {
		return err
	}

	return nil
}

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
