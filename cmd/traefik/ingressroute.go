// Package traefik
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package traefik

import (
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/kubectl"
)

var ingressRoute = &cobra.Command{
	Use:   "ingressroute",
	Short: "ingressroute Commands",
	Long:  `This command manages ingressroute.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := addIngressRoute()
		if err != nil {
			panic(err)
		}
	},
}

var ingressRouteTmpl = `---

apiVersion: traefik.containo.us/v1alpha1
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

  {{ if .Tls }}
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
	ServicePort int
}

type traefikIngressRouteTls struct {
	CertResolver string
}

func addIngressRoute() error {
	ingress := &traefikIngressRoute{
		Name:      "test",
		Namespace: "default",
		EntryPoints: []string{
			"websecure",
		},
		Rules: []traefikIngressRouteRule{
			{
				Host: "test.com",
				Services: []traefikIngressRouteRuleService{
					{
						ServiceName: "test",
						ServicePort: 80,
					},
				},
			},
		},
		Tls: &traefikIngressRouteTls{
			CertResolver: "default",
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
	Cmd.AddCommand(ingressRoute)
}
