package usecase

import (
	"context"
	"fmt"

	"github.com/snakesneaks/k8s-ssh4csh/pkg/pod"
	corev1 "k8s.io/api/core/v1"
)

type podUsecase struct {
	ctx        context.Context
	podHandler pod.CurrentPodHandler
}

type PodUsecase interface {
	GetCurrentPod() (*corev1.Pod, error)
	GetContainerInCurrentPod(targetContainerName string) (*corev1.Container, error)
}

var _ PodUsecase = &podUsecase{}

func NewPodUsecase(
	ctx context.Context,
	k8sUsecase K8sInClusterUsecase,
) (PodUsecase, error) {
	podHandler, err := pod.NewCurrentPodHandler(k8sUsecase.GetClientSet())
	if err != nil {
		return nil, err
	}
	podUsecase := &podUsecase{ctx, podHandler}
	return podUsecase, nil
}

func (u *podUsecase) GetContainerInCurrentPod(targetContainerName string) (*corev1.Container, error) {
	pod, err := u.podHandler.GetCurrentPod(u.ctx)
	if err != nil {
		return nil, err
	}

	var targetContainer *corev1.Container = nil
	for _, c := range pod.Spec.Containers {
		if c.Name == targetContainerName {
			targetContainer = &c
		}
	}
	if targetContainer == nil {
		names := make([]string, len(pod.Spec.Containers))
		for i, c := range pod.Spec.Containers {
			names[i] = c.Name
		}
		return nil, fmt.Errorf("target container %s not found. Container names in this pod are only the following: %v", targetContainer, names)
	}

	return targetContainer, nil
}

func (u *podUsecase) GetCurrentPod() (*corev1.Pod, error) {
	return u.podHandler.GetCurrentPod(u.ctx)
}
