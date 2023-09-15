package traefik

import (
	"context"
	"fmt"
	"github.com/zcubbs/zrun/cmd/vault"
	"github.com/zcubbs/zrun/internal/configs"
	"github.com/zcubbs/zrun/pkg/kubernetes"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

type DnsProvider string

const (
	Cloudflare DnsProvider = "cloudflare"
	OVH        DnsProvider = "ovh"
	Azure      DnsProvider = "azure"
)

func configureDNSChallenge() ([]string, error) {
	if dnsProviderString == "" {
		return nil, fmt.Errorf("dns provider is required")
	}

	if dnsProviderString == string(Cloudflare) {
		return configureCloudflare()
	}

	if dnsProviderString == string(OVH) {
		return configureOVH()
	}

	if dnsProviderString == string(Azure) {
		return configureAzure()
	}

	return nil, fmt.Errorf("unknown dns provider: %s", dnsProviderString)
}

func configureCloudflare() ([]string, error) {
	return nil, fmt.Errorf("cloudflare provider not implemented")
}

func configureOVH() ([]string, error) {
	if useVault {
		if err := configureOVHVault(); err != nil {
			return nil, err
		}
	} else {
		if err := configureOVHEnv(); err != nil {
			return nil, err
		}
	}

	return []string{
		fmt.Sprintf("--certificatesresolvers.%s.acme.dnschallenge.provider=ovh", dnsResolver),
	}, nil
}

func configureOVHEnv() error {
	// load env vars
	ovhEndpoint = os.Getenv(ovhEndpointEnvKey)
	ovhAppKey = os.Getenv(ovhAppKeyEnvKey)
	ovhAppSecret = os.Getenv(ovhAppSecretEnvKey)
	ovhConsumerKey = os.Getenv(ovhConsumerKeyEnvKey)

	// validate env vars
	if ovhEndpoint == "" {
		return fmt.Errorf("ovh endpoint is required")
	}

	if ovhAppKey == "" {
		return fmt.Errorf("ovh app key is required")
	}

	if ovhAppSecret == "" {
		return fmt.Errorf("ovh app secret is required")
	}

	if ovhConsumerKey == "" {
		return fmt.Errorf("ovh consumer key is required")
	}

	return nil
}

func configureOVHVault() error {
	// load vault secrets
	endpoint, err := vault.GetSecret(ovhEndpointVaultKey)
	if err != nil {
		return fmt.Errorf("failed to get ovh endpoint: %w", err)
	}
	appKey, err := vault.GetSecret(ovhAppKeyVaultKey)
	if err != nil {
		return fmt.Errorf("failed to get ovh app key: %w", err)
	}
	appSecret, err := vault.GetSecret(ovhAppSecretVaultKey)
	if err != nil {
		return fmt.Errorf("failed to get ovh app secret: %w", err)
	}
	consumerKey, err := vault.GetSecret(ovhConsumerKeyVaultKey)
	if err != nil {
		return fmt.Errorf("failed to get ovh consumer key: %w", err)
	}

	// validate vault secrets
	if endpoint == "" || appKey == "" || appSecret == "" || consumerKey == "" {
		return fmt.Errorf("failed to get ovh credentials from vault")
	}

	// set vault secrets
	ovhEndpoint = endpoint
	ovhAppKey = appKey
	ovhAppSecret = appSecret
	ovhConsumerKey = consumerKey

	return nil
}

func configureAzure() ([]string, error) {
	if useVault {
		if err := configureAzureVault(); err != nil {
			return nil, err
		}
	} else {
		if err := configureAzureEnv(); err != nil {
			return nil, err
		}
	}

	return []string{
		fmt.Sprintf("--certificatesresolvers.%s.acme.dnschallenge=true", dnsResolver),
		fmt.Sprintf("--certificatesresolvers.%s.acme.dnschallenge.provider=azure", dnsResolver),
		fmt.Sprintf("--certificatesresolvers.%s.acme.dnschallenge.azure.clientid=%s", dnsResolver, azureClientID),
		fmt.Sprintf("--certificatesresolvers.%s.acme.dnschallenge.azure.clientsecret=%s", dnsResolver, azureClientSecret),
	}, nil
}

func configureAzureEnv() error {
	// load env vars
	azureClientID = os.Getenv(azureClientIDEnvKey)
	azureClientSecret = os.Getenv(azureClientSecretEnvKey)

	// validate env vars
	if azureClientID == "" {
		return fmt.Errorf("azure client id is required")
	}

	if azureClientSecret == "" {
		return fmt.Errorf("azure client secret is required")
	}

	return nil
}

func configureAzureVault() error {
	// load vault secrets
	clientID, err := vault.GetSecret(azureClientIDVaultKey)
	if err != nil {
		return fmt.Errorf("failed to get azure client id: %w", err)
	}
	clientSecret, err := vault.GetSecret(azureClientSecretVaultKey)
	if err != nil {
		return fmt.Errorf("failed to get azure client secret: %w", err)
	}

	// validate vault secrets
	if clientID == "" || clientSecret == "" {
		return fmt.Errorf("failed to get azure credentials from vault")
	}

	// set vault secrets
	azureClientID = clientID
	azureClientSecret = clientSecret

	return nil
}

func createDnsSecret() error {
	if dnsProviderString == "" {
		return fmt.Errorf("dns provider is required")
	}

	if dnsProviderString == string(Cloudflare) {
		return createCloudflareDnsSecret()
	}

	if dnsProviderString == string(OVH) {
		return createOVHDnsSecret()
	}

	if dnsProviderString == string(Azure) {
		return createAzureDnsSecret()
	}

	return fmt.Errorf("unknown dns provider: %s", dnsProviderString)
}

func createCloudflareDnsSecret() error {
	return fmt.Errorf("cloudflare provider not implemented")
}

func createOVHDnsSecret() error {
	kubeconfig := configs.Config.Kubeconfig.Path
	err := kubernetes.CreateGenericSecret(
		context.Background(),
		kubeconfig,
		v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: "traefik-dns-account-credentials",
			},
			Data: map[string][]byte{
				"OVH_ENDPOINT":           []byte(ovhEndpoint),
				"OVH_APPLICATION_KEY":    []byte(ovhAppKey),
				"OVH_APPLICATION_SECRET": []byte(ovhAppSecret),
				"OVH_CONSUMER_KEY":       []byte(ovhConsumerKey),
				"TZ":                     []byte(dnsTz),
			},
		},
		[]string{traefikNamespace},
		true,
		false,
	)
	if err != nil {
		return err
	}
	return nil
}

func createAzureDnsSecret() error {

	return nil
}
