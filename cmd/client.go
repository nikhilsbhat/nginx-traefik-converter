package cmd

import (
	"github.com/nikhilsbhat/nginx-traefik-converter/pkg/log"
	"github.com/spf13/cobra"
)

func setCLIClient(cmd *cobra.Command, _ []string) error {
	logger = log.SetLogger(cliCfg.LogLevel)

	kubeConfig.SetLogger(logger)

	if cmd.Name() != "supported-annotations" {
		if err := kubeConfig.SetKubeClient(); err != nil {
			return err
		}
	}

	kubeConfig.SetKubeNameSpace()

	return nil
}
