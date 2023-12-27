package pod

// PRI is "pod resource identifier"
type PRI struct {
	podName, namespace string
}

func NewPRI(podName, namespace string) PRI {
	return PRI{
		podName:   podName,
		namespace: namespace,
	}
}

// CRI is "container resource identifier"
type CRI struct {
	containerName, podName, namespace string
}

func NewCRI(containerName, podName, namespace string) CRI {
	return CRI{
		containerName: containerName,
		podName:       podName,
		namespace:     namespace,
	}
}
