// Package kubectl
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package kubectl

import (
	"bytes"
	"context"
	"fmt"
	"github.com/zcubbs/zrun/bash"
	apiv1 "k8s.io/api/core/v1"
	errosv1 "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"text/template"
	"time"
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

func ApplyManifest(manifestTmpl string, data interface{}, debug bool) error {
	b, err := ApplyTmpl(manifestTmpl, data)
	if err != nil {
		return fmt.Errorf("failed to apply template \n %v", err)
	}

	// generate tmp file name
	fn := fmt.Sprintf("/tmp/tmpManifest_%s.yaml",
		time.Unix(time.Now().Unix(), 0).Format("20060102150405"),
	)

	if debug {
		// write tmp manifest
		err = os.WriteFile(fn, b, 0644)
		if err != nil {
			return fmt.Errorf("failed to write tmp manifest \n %v", err)
		}
	}

	err = bash.ExecuteCmd("kubectl", debug, "apply", "-f", fn)
	if err != nil {
		return fmt.Errorf("failed to apply manifest \n %s", err)
	}

	// delete tmp manifest
	err = os.Remove(fn)
	if err != nil {
		return fmt.Errorf("failed to delete tmp manifest \n %v", err)
	}
	return nil
}

func ApplyTmpl(tmplStr string, tmplData interface{}) ([]byte, error) {
	tmpl, err := template.New("tmpManifest").Parse(tmplStr)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, tmplData); err != nil {
		return nil, err
	}

	fmt.Println(buf.String())

	return buf.Bytes(), nil
}
