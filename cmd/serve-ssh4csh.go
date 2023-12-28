package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/snakesneaks/k8s-ssh4csh/internal/config"
	"github.com/snakesneaks/k8s-ssh4csh/internal/usecase"
	"github.com/snakesneaks/k8s-ssh4csh/pkg/pod"
	"github.com/snakesneaks/k8s-ssh4csh/pkg/ssh"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
// it serves SSH connection and execute command
var serveCmd = &cobra.Command{
	Use:          "serve-ssh4csh",
	Short:        "server ssh connection",
	Long:         `using this, users can connect to containers in the same pod.`,
	RunE:         serve,
	SilenceUsage: true,
}

var addr string

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&addr, "address", "a", ":2222", "Address to listen")
}

func serve(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	k8sInClusterUsecase, err := usecase.NewK8sInClusterUsecase()
	if err != nil {
		return err
	}

	env, err := config.LoadSsh4CshServerEnv()
	if err != nil {
		return err
	}

	publicKey, err := ssh.ParsePubKeyString(env.CONFIG.SSH_PUB_KEY)
	if err != nil {
		return fmt.Errorf(" parse ssh public key \"%s\" : %v", env.CONFIG.SSH_PUB_KEY, err)
	}

	sshUsecase, err := usecase.NewSshUsecase(
		env.CONFIG.SERVER_PORT,
		ctx,
		publicKey,
		k8sInClusterUsecase)
	if err != nil {
		return err
	}

	podUsecase, err := usecase.NewPodUsecase(ctx, k8sInClusterUsecase)
	if err != nil {
		return err
	}
	currentPod, err := podUsecase.GetCurrentPod()
	if err != nil {
		return err
	}
	/*
		if _, err := podUsecase.GetContainerInCurrentPod(env.TARGET_CONTAINER); err != nil {
			return err
		}
	*/
	cri := pod.NewCRI(env.TARGET_CONTAINER, currentPod.GetName(), currentPod.GetNamespace())

	log.Printf("port: %s\nssh-pub-key: %s\n", env.CONFIG.SERVER_PORT, env.CONFIG.SSH_PUB_KEY)
	log.Printf("targetContainer: %s", env.TARGET_CONTAINER)
	log.Println("server started!!")

	var serverErr error
	go func() error {
		if err := sshUsecase.StartSshServer4ContainerShell(cri, env.CONFIG.SHELL); err != nil {
			serverErr = err
			cancel()
			return err
		}
		return nil
	}()

	<-ctx.Done()

	if serverErr != nil {
		return serverErr
	}

	if err := sshUsecase.Shutdown(); err != nil {
		log.Println("Server shutdown!!")
		return err
	}

	return nil
}
