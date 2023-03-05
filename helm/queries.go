package helm

import (
	"github.com/zcubbs/zrun/defaults"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/release"
	"log"
)

func GetAllReleases(kubeconfig string) ([]*release.Release, error) {
	if kubeconfig == "" {
		kubeconfig = defaults.DefaultKubeConfigPath
	}

	actionConfig := new(action.Configuration)
	err := actionConfig.Init(kube.GetConfig(kubeconfig, "", ""), "", "", log.Printf)
	if err != nil {
		return nil, err
	}

	_releases, err := action.NewList(actionConfig).Run()
	if err != nil {
		return nil, err
	}

	return _releases, nil
}
