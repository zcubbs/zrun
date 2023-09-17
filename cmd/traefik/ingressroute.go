// Package traefik
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package traefik

import (
	"fmt"
	"github.com/spf13/cobra"
	kubectl "github.com/zcubbs/x/kubernetes"
	"github.com/zcubbs/x/must"
	"github.com/zcubbs/x/progress"
	"github.com/zcubbs/x/style"
)

var (
	ingressRouteName        string   // ingressroute name
	ingressRouteNamespace   string   // ingressroute namespace
	ingressRouteEntryPoints []string // ingressroute entrypoints
	ingressRoutePath        string   // ingressroute path
	ingressRouteService     string   // ingressroute service
	ingressRoutePort        string   // ingressroute port
	ingressRouteTls         bool     // ingressroute tls
	ingressRouteTlsSecret   string   // ingressroute tls secret
)

var ingressRoute = &cobra.Command{
	Use:   "ingressroute",
	Short: "ingressroute Commands",
	Long:  `This command manages ingressroute.`,
	Run: func(cmd *cobra.Command, args []string) {
		style.PrintColoredHeader("add traefik ingress-route")
		must.Succeed(
			progress.RunTask(func() error {
				err := addIngressRoute()
				if err != nil {
					return err
				}
				return nil
			}, true))
	},
}

var ingressRouteTmpl = `---

apiVersion: traefik.io/v1alpha1
kind: IngressRoute
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  entryPoints:
	{{- range .EntryPoints }}
    - websecure
	{{- end }}
  routes:
  {{- range .Rules }}
    - kind: Rule
      match: Host(` + "`{{ .Host }}`" + `)
      {{- range .Services }}
      services:
        - name: {{ .ServiceName }}
          port: {{ .ServicePort }}
      {{- end }}
  {{- end }}

  {{ if .Tls.IsEnabled }}
  tls:
	certResolver: {{ .Tls.CertResolver }}
  {{- end }}
`

type traefikIngressRoute struct {
	Name        string
	Namespace   string
	EntryPoints []string
	Rules       []traefikIngressRouteRule
	Tls         *traefikIngressRouteTls
}

type traefikIngressRouteRule struct {
	Host     string
	Services []traefikIngressRouteRuleService
}

type traefikIngressRouteRuleService struct {
	ServiceName string
	ServicePort string
}

type traefikIngressRouteTls struct {
	IsEnabled    bool
	CertResolver string
}

func addIngressRoute() error {
	ingress := &traefikIngressRoute{
		Name:        ingressRouteName,
		Namespace:   ingressRouteNamespace,
		EntryPoints: ingressRouteEntryPoints,
		Rules: []traefikIngressRouteRule{
			{
				Host: ingressRoutePath,
				Services: []traefikIngressRouteRuleService{
					{
						ServiceName: ingressRouteService,
						ServicePort: ingressRoutePort,
					},
				},
			},
		},
		Tls: &traefikIngressRouteTls{
			IsEnabled:    ingressRouteTls,
			CertResolver: ingressRouteTlsSecret,
		},
	}

	// Apply template
	verbose := Cmd.Flag("verbose").Value.String() == "true"
	err := kubectl.ApplyManifest(ingressRouteTmpl, ingress, verbose)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	ingressRoute.Flags().StringVarP(&ingressRouteName, "name", "n", "", "ingressroute name")
	ingressRoute.Flags().StringVarP(&ingressRouteNamespace, "namespace", "s", "", "ingressroute namespace")
	ingressRoute.Flags().StringSliceVarP(&ingressRouteEntryPoints, "entrypoints", "e", []string{}, "ingressroute entrypoints")
	ingressRoute.Flags().StringVarP(&ingressRoutePath, "path", "p", "", "ingressroute path")
	ingressRoute.Flags().StringVarP(&ingressRouteService, "service", "S", "", "ingressroute service")
	ingressRoute.Flags().StringVarP(&ingressRoutePort, "port", "P", "", "ingressroute port")
	ingressRoute.Flags().BoolVarP(&ingressRouteTls, "tls", "t", false, "ingressroute tls")
	ingressRoute.Flags().StringVarP(&ingressRouteTlsSecret, "tls-secret", "T", "", "ingressroute tls secret")

	// required flags
	err := ingressRoute.MarkFlagRequired("name")
	if err != nil {
		fmt.Println(err)
	}
	err = ingressRoute.MarkFlagRequired("namespace")
	if err != nil {
		fmt.Println(err)
	}
	err = ingressRoute.MarkFlagRequired("entrypoints")
	if err != nil {
		fmt.Println(err)
	}
	err = ingressRoute.MarkFlagRequired("path")
	if err != nil {
		fmt.Println(err)
	}
	err = ingressRoute.MarkFlagRequired("service")
	if err != nil {
		fmt.Println(err)
	}
	err = ingressRoute.MarkFlagRequired("port")
	if err != nil {
		fmt.Println(err)
	}

	Cmd.AddCommand(ingressRoute)
}
