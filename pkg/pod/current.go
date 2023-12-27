package pod

import (
	"context"
	"errors"
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type CurrentPodHandler interface {
	GetCurrentPodName() string
	GetCurrentNamespace() string
	GetCurrentPod(ctx context.Context) (*corev1.Pod, error)
}

type currentPodHandler struct {
	currentPodName string
	currentNS      string
	clientSet      *kubernetes.Clientset
}

var _ CurrentPodHandler = &currentPodHandler{}

func NewCurrentPodHandler(clientSet *kubernetes.Clientset) (CurrentPodHandler, error) {
	currentPodName, err := getCurrentPodName()
	if err != nil {
		return nil, err
	}
	currentNS, err := getCurrentNamespace()
	if err != nil {
		return nil, err
	}
	return &currentPodHandler{
		currentPodName: currentPodName,
		currentNS:      currentNS,
		clientSet:      clientSet,
	}, nil
}
func getCurrentPodName() (string, error) {
	var err error
	podName, ok := os.LookupEnv("HOSTNAME")
	if !ok {
		err = errors.New("failed to find environmental variable of HOSTNAME, which should be automatically assigned to the pod")
	}
	return podName, err
}

func getCurrentNamespace() (string, error) {
	namespaceFilePath := "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	data, err := os.ReadFile(namespaceFilePath)
	if err != nil {
		return "", err
	}
	if len(data) == 0 {
		return "", fmt.Errorf("namespace was not found in the file %s", namespaceFilePath)
	}
	return string(data), nil
}

func (h *currentPodHandler) GetCurrentPodName() string {
	return h.currentPodName
}

func (h *currentPodHandler) GetCurrentNamespace() string {
	return h.currentNS
}
func (h *currentPodHandler) GetCurrentPod(ctx context.Context) (*corev1.Pod, error) {
	currentPod, err := h.GetPod(ctx, h.currentNS, h.currentPodName)
	if err != nil {
		return nil, err
	}
	return currentPod, err
}

func (h *currentPodHandler) GetPod(ctx context.Context, ns, podName string) (*corev1.Pod, error) {
	p, err := h.clientSet.CoreV1().Pods(h.currentNS).Get(ctx, h.currentPodName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return p, err
}
