package k3s

import (
	"fmt"
	"github.com/zcubbs/zrun/bash"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"runtime"
)

func InstallK9s() error {
	fmt.Println("Installing k9s...")
	fmt.Printf("Arch: %s_%s\n", runtime.GOOS, runtime.GOARCH)
	// wget -O /tmp/k9s.tar.gz https://github.com/derailed/k9s/releases/download/v0.27.3/k9s_Linux_amd64.tar.gz
	err := bash.ExecuteCmd(
		"wget",
		"-O",
		"/tmp/k9s.tar.gz",
		fmt.Sprintf("https://github.com/derailed/k9s/releases/download/v0.27.3/k9s_%s_%s.tar.gz",
			cases.Title(language.Und).String(runtime.GOOS),
			runtime.GOARCH,
		),
	)
	if err != nil {
		return err
	}

	// tar -xvf k9s.tar.gz
	err = bash.ExtractTarGzWithFile("/tmp/k9s.tar.gz", "k9s", "/tmp")
	if err != nil {
		return err
	}

	// mv /tmp/k9s /usr/local/bin/k9s
	err = bash.ExecuteCmd(
		"mv",
		"/tmp/k9s",
		"/usr/local/bin/k9s",
	)
	if err != nil {
		return err
	}

	// chmod +x /usr/local/bin/k9s
	err = bash.Chmod("/usr/local/bin/k9s", 0700)
	if err != nil {
		return err
	}

	return nil
}
