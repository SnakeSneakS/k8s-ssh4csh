package pod

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gliderlabs/ssh"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type sshHandler struct {
	restConfig *rest.Config
	restClient *rest.RESTClient
}
type SshHandler interface {
	Ssh4ContainerShellHandler(cri CRI, cmd string) func(ssh.Session)
}

func NewSshHandler(
	config *rest.Config,
) (SshHandler, error) {
	restConfig := rest.CopyConfig(config)
	if err := setK8sRestConfigDefault(restConfig); err != nil {
		return nil, err
	}
	restClient, err := rest.RESTClientFor(restConfig)
	if err != nil {
		return nil, err
	}
	return &sshHandler{
		restConfig: restConfig,
		restClient: restClient,
	}, nil
}

type User struct {
	PublicKey ssh.PublicKey
	User      string
	Pod       string
	Namespace string
}

var ErrDestination = errors.New("can't find destination")

const AuthorizedKeyAnnotation = "ssh.barpilot.io/publickey"
const CommandAnnotation = "ssh.barpilot.io/command"
const PrefixCommandAnnotation = "ssh.barpilot.io/prefix-command"

func (h *sshHandler) Ssh4ContainerShellHandler(cri CRI, cmd string) func(ssh.Session) {
	return func(s ssh.Session) {
		ctx := s.Context()
		cmds := s.Command()
		cmds = append([]string{cmd}, cmds...)
		//shlex.Split()
		_, cWindows, hasPTY := s.Pty()
		queue := sizeQueue{C: cWindows}

		req := h.restClient.Post().
			Resource("pods").
			Name(cri.podName).
			Namespace(cri.namespace).
			SubResource("exec")
		req.VersionedParams(&v1.PodExecOptions{
			Command:   cmds,
			Container: cri.containerName,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       hasPTY,
		}, clientgoscheme.ParameterCodec)

		log.Printf("exec generated URL: %s", req.URL())
		executer, err := remotecommand.NewSPDYExecutor(h.restConfig, "POST", req.URL())
		if err != nil {
			log.Printf("fail create NewSPDYExecutor for url '%s': %v", req.URL(), err)
			s.Stderr().Write([]byte(ErrDestination.Error()))
			s.Exit(1)
			return
		}

		// if !ok {
		if err := executer.StreamWithContext(ctx, remotecommand.StreamOptions{
			Tty:               hasPTY,
			Stdin:             s,
			Stdout:            s,
			Stderr:            s.Stderr(),
			TerminalSizeQueue: queue,
		}); err != nil {
			log.Printf("fail to exec Stream: %v", err)
			s.Stderr().Write([]byte(ErrDestination.Error()))
			s.Exit(1)
			return
		}

		return
	}
}

type sizeQueue struct {
	C <-chan ssh.Window
}

func (s sizeQueue) Next() *remotecommand.TerminalSize {
	size, ok := <-s.C
	if !ok {
		return nil
	}
	tSize := &remotecommand.TerminalSize{
		Width:  uint16(size.Width),
		Height: uint16(size.Height),
	}

	return tSize
}

// getCRIFromSshUser sshUser should be in the format of CONTIANER.POD.NAMESPACE@DOMAIN
func getCRIFromSshUser(sshUser string) (CRI, error) {
	user, _, ok := strings.Cut(sshUser, "@")
	if !ok {
		return CRI{}, errors.New("can't parse ssh user")
	}

	ul := strings.Split(user, ".")
	if len(ul) != 3 {
		return CRI{}, fmt.Errorf("sshUser should be in the format of CONTAINER.POD.NAMESPACE@DOMAIN")
	}
	cri := CRI{
		containerName: ul[0],
		podName:       ul[1],
		namespace:     ul[2],
	}
	return cri, nil
}

func setK8sRestConfigDefault(config *rest.Config) error {
	config.GroupVersion = &schema.GroupVersion{Group: "", Version: "v1"}

	if config.APIPath == "" {
		config.APIPath = "/api"
	}

	if config.NegotiatedSerializer == nil {
		config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	}

	return rest.SetKubernetesDefaults(config)
}
