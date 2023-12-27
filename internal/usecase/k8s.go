package usecase

import (
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type k8sInClusterUsecase struct {
	config    *rest.Config
	clientSet *kubernetes.Clientset
}

type K8sInClusterUsecase interface {
	GetConfig() *rest.Config
	GetClientSet() *kubernetes.Clientset
}

var _ K8sInClusterUsecase = &k8sInClusterUsecase{}

func NewK8sInClusterUsecase() (K8sInClusterUsecase, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &k8sInClusterUsecase{
		config,
		clientSet,
	}, nil
}

func (u k8sInClusterUsecase) GetClientSet() *kubernetes.Clientset {
	return u.clientSet
}
func (u k8sInClusterUsecase) GetConfig() *rest.Config {
	return u.config
}
