package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cockroachdb/errors"
	"github.com/snakesneaks/k8s-ssh4csh/internal/usecase"

	"github.com/spf13/cobra"
)

// echoContainersCmd represents the echoContainers command
// it displays containers
var echoContainersCmd = &cobra.Command{
	Use:          "echo-containers",
	Short:        "echo containers in this pod",
	Long:         `using this, users can check containers in this pod.`,
	RunE:         echoContainers,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(echoContainersCmd)
}

func echoContainers(cmd *cobra.Command, args []string) error {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	k8sInClusterUsecase, err := usecase.NewK8sInClusterUsecase()
	if err != nil {
		return errors.WithStack(err)
	}

	podUsecase, err := usecase.NewPodUsecase(ctx, k8sInClusterUsecase)
	if err != nil {
		return errors.WithStack(err)
	}

	currentPod, err := podUsecase.GetCurrentPod()
	if err != nil {
		return errors.WithStack(err)
	}

	containerNames := make([]string, len(currentPod.Spec.Containers))
	for i, c := range currentPod.Spec.Containers {
		containerNames[i] = c.Name
	}
	log.Printf("container names in this pod (%s): %#v", currentPod.Name, containerNames)
	return nil
}
