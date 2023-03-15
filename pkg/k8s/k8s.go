package k8s

import (
	"context"

	"github.com/bitnami-labs/sealed-secrets/pkg/apis/sealedsecrets/v1alpha1"
	ssscheme "github.com/bitnami-labs/sealed-secrets/pkg/client/clientset/versioned/scheme"
	"github.com/pkg/errors"
	"github.com/takutakahashi/k2ssm/pkg/output"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/strings/slices"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Kubernetes struct {
	c client.Client
}

type GatherSecretsOpt struct {
	Namespace        string
	MatchLabels      map[string]string
	FromSealedSecret bool
}

func NewKubernetes() (*Kubernetes, error) {
	scheme := runtime.NewScheme()
	clientgoscheme.AddToScheme(scheme)
	ssscheme.AddToScheme(scheme)
	k8sClient, err := client.New(ctrl.GetConfigOrDie(), client.Options{Scheme: scheme})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create k8s client")
	}

	return &Kubernetes{c: k8sClient}, nil
}
func (k *Kubernetes) GatherSecrets(ctx context.Context, opt GatherSecretsOpt) (output.Output, error) {
	ret := output.Output{Secrets: []output.OutputSecret{}}
	targetNames := []string{}
	if opt.FromSealedSecret || true {
		ssl := &v1alpha1.SealedSecretList{}
		if err := k.c.List(ctx, ssl, &client.ListOptions{
			Namespace:     opt.Namespace,
			LabelSelector: labels.Set(opt.MatchLabels).AsSelector(),
		}); err != nil {
			return ret, errors.Wrap(err, "failed to list secrets")
		}
		for _, ss := range ssl.Items {
			targetNames = append(targetNames, ss.Name)
		}
	}
	sl := &corev1.SecretList{}
	if err := k.c.List(ctx, sl, &client.ListOptions{
		Namespace:     opt.Namespace,
		LabelSelector: labels.Set(opt.MatchLabels).AsSelector(),
	}); err != nil {
		return ret, errors.Wrap(err, "failed to list secrets")
	}
	res := []corev1.Secret{}
	if len(targetNames) != 0 {

		for _, s := range sl.Items {
			if slices.Contains(targetNames, s.Name) {
				res = append(res, s)
			}
		}
	} else {
		res = sl.Items
	}
	for _, s := range res {
		d := map[string]string{}
		for k, v := range s.Data {
			d[k] = string(v)
		}
		ret.Secrets = append(ret.Secrets, output.OutputSecret{
			Name:      s.Name,
			Namespace: s.Namespace,
			Type:      string(s.Type),
			Data:      d,
		})
	}
	return ret, nil
}
