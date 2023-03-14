package k8s

import (
	ssscheme "github.com/bitnami-labs/sealed-secrets/pkg/client/clientset/versioned/scheme"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Kubernetes struct {
	c client.Client
}

func NewKubernetes() (*Kubernetes, error) {
	scheme := runtime.NewScheme()
	clientgoscheme.AddToScheme(scheme)
	ssscheme.AddToScheme(scheme)
	k8sClient, err := client.New(ctrl.GetConfigOrDie(), client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}

	return &Kubernetes{c: k8sClient}, nil
}
