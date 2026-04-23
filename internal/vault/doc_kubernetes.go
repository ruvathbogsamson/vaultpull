// Package vault provides Kubernetes auth method support for HashiCorp Vault.
//
// The KubernetesClient authenticates a service account against Vault using the
// Kubernetes auth method. It exchanges a Kubernetes service account JWT and a
// configured role name for a Vault client token.
//
// # Usage
//
//	client, err := vault.NewKubernetesClient("http://vault:8200", "kubernetes")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	jwt, _ := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
//	resp, err := client.Login("my-app-role", string(jwt))
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Vault token:", resp.ClientToken)
//
// The mount path defaults to "kubernetes" if not specified.
package vault
