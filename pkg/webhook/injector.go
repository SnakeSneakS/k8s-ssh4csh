package webhook

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var ssh4cshDefaultServiceAccount = "ssh4csh-service-account"
var ssh4cshContainerDefault = corev1.Container{
	Name:  "ssh4csh-sidecar",
	Image: "",
	Resources: corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("0.5"),
			corev1.ResourceMemory: resource.MustParse("100Mi"),
		},
	},
}

func Ssh4CshSidecarInjector(ctx context.Context, r admission.Request) admission.Response {
	switch r.RequestKind.Kind {
	case "Pod":
		var pod corev1.Pod
		decoder := admission.NewDecoder(clientgoscheme.Scheme)
		if err := decoder.Decode(r, &pod); err != nil {
			return admission.Denied(fmt.Sprintf("failed to decode request: %v", err))
		}
		pod.Spec.ServiceAccountName = ssh4cshDefaultServiceAccount
		pod.Spec.Containers = append(pod.Spec.Containers, ssh4cshContainerDefault)
		break
	default:
		return admission.Denied("not implemented")
	}

	return admission.Allowed("allowed")
}
