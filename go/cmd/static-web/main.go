package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"go.f110.dev/mono/go/pkg/http/httpserver"
	"go.f110.dev/mono/go/pkg/logger"
)

type Mode string

const (
	ModeSPA    Mode = "spa"
	ModeSimple Mode = "simple"
)

func staticWeb() error {
	var documentRoot, listenAddr, mode string
	cmd := &cobra.Command{
		Use:   "static-web",
		Short: "Serve static files",
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return logger.Init()
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			http.Handle("/favicon.ico", http.NotFoundHandler())
			switch Mode(mode) {
			case ModeSPA:
				http.Handle("/", httpserver.SinglePageApplication(documentRoot))
			case ModeSimple, "":
				http.Handle("/", http.FileServer(http.Dir(documentRoot)))
			}

			s := &http.Server{
				Addr:    listenAddr,
				Handler: http.DefaultServeMux,
			}
			go func() {
				ctx := cmd.Context()
				<-ctx.Done()
				logger.Log.Info("Shutdown")
				s.Shutdown(context.Background())
			}()
			logger.Log.Info("Start server", zap.String("addr", listenAddr), zap.String("root", documentRoot), zap.String("mode", mode))
			if err := s.ListenAndServe(); err == http.ErrServerClosed {
				return nil
			} else if err != nil {
				return err
			}
			return nil
		},
	}
	logger.Flags(cmd.Flags())
	cmd.Flags().StringVar(&documentRoot, "document-root", "", "The document root")
	cmd.Flags().StringVar(&listenAddr, "listen-addr", ":8050", "Listen address")
	cmd.Flags().StringVar(&mode, "mode", string(ModeSimple), "")

	ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancelFunc()
	return cmd.ExecuteContext(ctx)
}

func main() {
	if err := staticWeb(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}