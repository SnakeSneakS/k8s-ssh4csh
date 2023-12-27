package usecase

import (
	"context"
	"errors"
	"fmt"

	gossh "golang.org/x/crypto/ssh"

	"github.com/gliderlabs/ssh"
	"github.com/snakesneaks/k8s-ssh-sidecar/pkg/pod"
	sshserver "github.com/snakesneaks/k8s-ssh-sidecar/pkg/ssh"
)

type sshUsecase struct {
	ctx        context.Context
	port       string
	publicKey  ssh.PublicKey
	server     *ssh.Server
	sshHandler pod.SshHandler
	podHandler pod.CurrentPodHandler
}

type SshUsecase interface {
	StartSshServer4ContainerShell(cri pod.CRI, cmd string) error
	Shutdown() error
}

var _ SshUsecase = &sshUsecase{}

func NewSshUsecase(
	port string,
	ctx context.Context,
	publicKey gossh.PublicKey,
	k8sUsecase K8sInClusterUsecase,
) (SshUsecase, error) {
	podHandler, err := pod.NewCurrentPodHandler(k8sUsecase.GetClientSet())
	if err != nil {
		return nil, err
	}
	sshHandler, err := pod.NewSshHandler(k8sUsecase.GetConfig())
	if err != nil {
		return nil, err
	}
	sshUsecase := &sshUsecase{
		port:       port,
		ctx:        ctx,
		publicKey:  publicKey,
		sshHandler: sshHandler,
		podHandler: podHandler,
	}
	return sshUsecase, nil
}

// StartSshServer4ContainerShell runs ssh server which is connected to container shell
func (s *sshUsecase) StartSshServer4ContainerShell(cri pod.CRI, cmd string) error {
	addr := fmt.Sprintf(":%s", s.port)

	handler := s.sshHandler.Ssh4ContainerShellHandler(cri, cmd) //ssh.DefaultHandler
	rl := sshserver.NewRateLimiter()

	server := ssh.Server{
		Addr:              addr,
		Handler:           handler,                                 //pod.SshHandler(clientset, config),
		PublicKeyHandler:  sshserver.PublicKeyHandler(s.publicKey), //PublicKeyHandler(clientset),
		PasswordHandler:   func(ctx ssh.Context, password string) bool { return false },
		SubsystemHandlers: map[string]ssh.SubsystemHandler{
			//"sftp": SftpHandler(),
		},
		ConnCallback:             rl.ConnCallback(),
		ConnectionFailedCallback: rl.ConnectionFailedCallback(),
	}

	/*
		var lc net.ListenConfig
		l, err := lc.Listen(s.ctx, "tcp", fmt.Sprintf(":%s", s.port))
		if err != nil {
			return err
		}

		if err := server.Serve(l); err != nil {
			return err
		}
	*/
	s.server = &server
	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (s *sshUsecase) Shutdown() error {
	if s.server == nil {
		return errors.New("server is not started or set")
	}

	if err := s.server.Shutdown(s.ctx); err != nil {
		return err
	}
	return nil
}
