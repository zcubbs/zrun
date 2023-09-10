package kubernetes

import (
	"context"
	"encoding/base64"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type ContainerRegistrySecret struct {
	Name     string
	Server   string
	Username string
	Password string
	Email    string
}

func CreateContainerRegistrySecret(
	ctx context.Context,
	kubeconfig string,
	secretConfig ContainerRegistrySecret,
	namespaces []string,
	replace bool,
	debug bool) error {
	auth := fmt.Sprintf("%s:%s", secretConfig.Username, secretConfig.Password)
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))

	data := map[string][]byte{
		".dockerconfigjson": []byte(fmt.Sprintf(`{
			"auths": {
				"%s": {
					"username": "%s",
					"password": "%s",
					"email": "%s",
					"auth": "%s"
				}
			}
		}`, secretConfig.Server, secretConfig.Username,
			secretConfig.Password, secretConfig.Email, encodedAuth)),
	}

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secretConfig.Name,
		},
		Data: data,
		Type: v1.SecretTypeDockerConfigJson,
	}

	cs := GetClientSet(kubeconfig)

	for _, namespace := range namespaces {
		created, err := cs.CoreV1().Secrets(namespace).Create(ctx, secret, metav1.CreateOptions{})
		if err != nil {
			if strings.Contains(err.Error(), "already exists") && replace {
				_, err := cs.CoreV1().Secrets(namespace).Update(ctx, secret, metav1.UpdateOptions{})
				return err
			}
			return fmt.Errorf("failed to create secret: %v", err)
		}

		if debug {
			fmt.Printf("Created secret %s\n", created.String())
		}
	}

	return nil
}
