package ci

import (
	"flag"
	"fmt"
	"github.com/zcubbs/dagger-utils/container"
	"github.com/zcubbs/dagger-utils/golang"
	"github.com/zcubbs/dagger-utils/types"
)

var (
	binName, onlyBuild, builderImage, registryServer, registryUsername, registryPassword, registryEmail, registryImageName, registryImageTag string
)

func init() {
	flag.StringVar(&binName, "bin-name", "zrun", "registry server")
	flag.StringVar(&binName, "n", "zrun", "registry server")

	flag.StringVar(&builderImage, "builder-image", "paketobuildpacks/builder:base", "builder image to use")
	flag.StringVar(&builderImage, "b", "paketobuildpacks/builder:base", "builder image to use")

	flag.StringVar(&registryServer, "registry-server", "ttl.sh", "registry server")
	flag.StringVar(&registryServer, "s", "ttl.sh", "registry server")

	flag.StringVar(&registryUsername, "registry-username", "", "registry username")
	flag.StringVar(&registryUsername, "u", "", "registry username")

	flag.StringVar(&registryPassword, "registry-password", "", "registry password")
	flag.StringVar(&registryPassword, "p", "", "registry password")

	flag.StringVar(&registryEmail, "registry-email", "", "registry email")
	flag.StringVar(&registryEmail, "e", "", "registry email")

	flag.StringVar(&registryImageName, "registry-image-name", "", "registry image name")
	flag.StringVar(&registryImageName, "i", "", "registry image name")

	flag.StringVar(&registryImageTag, "registry-image-tag", "latest", "registry image tag")
	flag.StringVar(&registryImageTag, "a", "latest", "registry image tag")

	flag.StringVar(&onlyBuild, "only-build", "false", "only build the image, do not generate SBOM/vuln report")
	flag.StringVar(&onlyBuild, "", "false", "only build the image")

}

func main() {
	flag.Parse()

	goBuilder := golang.Builder{
		GoOptions: types.GoOptions{
			BinName: binName,
		},
	}

	// GO LINT
	err := goBuilder.GoLint()
	if err != nil {
		panic(err)
	}

	// GO TEST
	err = goBuilder.GoTest()
	if err != nil {
		panic(err)
	}

	// GO BUILD
	err = goBuilder.GoBuild()
	if err != nil {
		panic(err)
	}

	// GO BUILD CONTAINER IMAGE
	builder := container.ImageBuilder{
		RegistryInfo: types.RegistryInfo{
			RegistryServer:   registryServer,
			RegistryUsername: registryUsername,
			RegistryPassword: registryPassword,
			RegistryEmail:    registryEmail,
		},
	}

	if registryImageName == "" {
		registryImageName = binName
	}

	imgName, err := builder.Build(registryImageName, registryImageTag)
	if err != nil {
		panic(err)
	}

	fmt.Println("built image: ", *imgName)

	// GO SCAN CONTAINER IMAGE
	scanner := container.Scanner{}
	err = scanner.GenerateVulnReport(
		fmt.Sprintf("%s/%s:%s", registryServer, registryImageName, registryImageTag),
	)
	if err != nil {
		panic(err)
	}

	vulns, err := scanner.ScanVuln()
	if err != nil {
		panic(err)
	}

	// Print out the number of vulnerabilities found for each severity level
	levels, fixes := scanner.ParseVulnForSeverityLevels(vulns)
	for level, count := range levels {
		fmt.Printf("Found %d %s vulnerabilities\n", count, level)
	}
	fmt.Printf("%d vulnerabilities have fixes available\n", fixes)
}
