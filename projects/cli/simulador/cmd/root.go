package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "simulador",
	Short: "CLI para gestão do ciclo de vida do simulador HubSaúde",
}
