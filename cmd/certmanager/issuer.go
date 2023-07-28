// Package certmanager
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package certmanager

import (
	"github.com/spf13/cobra"
	"github.com/zcubbs/zrun/kubectl"
	"log"
)

const letsEncryptIssuerTmpl = `---

apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: {{ .IssuerName }}
  namespace: kube-system
spec:
  acme:
    email: {{ .IssuerEmail }}
    server: {{ .IssuerServer }}
    privateKeySecretRef:
      name: issuer-account-key
    solvers:
      - http01:
          ingress:
            class: {{ .IngressClassResolver }}
`

type Issuer struct {
	IssuerName           string
	IssuerEmail          string
	IssuerServer         string
	IngressClassResolver string
}

var (
	issuerName           string
	issuerEmail          string
	issuerServer         string
	ingressClassResolver string
)

// issuer represents the list command
var issuer = &cobra.Command{
	Use:   "issuer",
	Short: "setup cert-manager issuer",
	Long:  `setup cert-manager issuer. Note: currently only supports letsencrypt`,
	Run: func(cmd *cobra.Command, args []string) {
		verbose := Cmd.Flag("verbose").Value.String() == "true"
		err := kubectl.ApplyManifest(letsEncryptIssuerTmpl, Issuer{
			IssuerName:           issuerName,
			IssuerEmail:          issuerEmail,
			IssuerServer:         issuerServer,
			IngressClassResolver: ingressClassResolver,
		}, verbose)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	// parse flags
	issuer.Flags().StringVar(&issuerName, "name", "", "issuer name")
	issuer.Flags().StringVar(&issuerEmail, "email", "", "issuer email")
	issuer.Flags().StringVar(&issuerServer, "server", "", "issuer server")
	issuer.Flags().StringVar(&ingressClassResolver, "ingress-class-resolver", "cert-manager-resolver", "ingress class resolver")

	if err := issuer.MarkFlagRequired("name"); err != nil {
		log.Println(err)
	}

	if err := issuer.MarkFlagRequired("email"); err != nil {
		log.Println(err)
	}

	if err := issuer.MarkFlagRequired("server"); err != nil {
		log.Println(err)
	}

	Cmd.AddCommand(issuer)
}
