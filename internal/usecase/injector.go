package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type injectorUsecase struct {
	ctx               context.Context
	port              int //host:port
	webhookInjectPath string
	server            *http.Server
}

type InjectorUsecase interface {
}

var _ InjectorUsecase = &injectorUsecase{}

func NewInjectorUsecase(
	ctx context.Context,
	port int,
	webhookInjectPath string,
) InjectorUsecase {
	return &injectorUsecase{
		ctx:               ctx,
		port:              port,
		webhookInjectPath: webhookInjectPath,
	}
}

// NewSsh4CshSidecarInjectorServer runs a server which use Mutating Admission Webhook
// ref: https://github.com/kubernetes-sigs/controller-runtime/blob/main/pkg/webhook/example_test.go#L111
// Note that this assumes and requires a valid TLS
// cert and key at the default locations
// tls.crt and tls.key.
func (u *injectorUsecase) NewSsh4CshSidecarInjectorServer() error {
	mux := http.NewServeMux()

	h, err := admission.StandaloneWebhook(&admission.Webhook{
		Handler: admission.HandlerFunc(func(ctx context.Context, r admission.Request) admission.Response {
			return admission.Allowed("allowed")
		}),
	}, admission.StandaloneOptions{
		MetricsPath: "/mutating",
	})
	if err != nil {
		return err
	}
	mux.Handle(u.webhookInjectPath, h)
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", u.port),
		Handler: mux,
	}
	u.server = &server

	// start webhook server in new rountine
	if err := server.ListenAndServe(); err != nil {
		return fmt.Errorf("Failed to listen and serve webhook server: %v", err)
	}
	return nil
}

func (u *injectorUsecase) ShutdownServer(ctx context.Context) error {
	if u.server == nil {
		return errors.New("server is not started or set")
	}
	if err := u.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
