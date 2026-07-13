package main

import (
	"fmt"
	"log"

	"github.com/ing-bank/kaet/internal/cmd"
	"github.com/ing-bank/kaet/internal/logger"
	"github.com/spf13/cobra"
	"k8s.io/client-go/util/homedir"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("failed to execute command: %s\n", err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "kaet",
	Short: "Kubernetes Auto Exploit Tool - KAET",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Configure logger
		logger.Configure(cmd)
	},
}

func init() {
	rootCmd.AddCommand(cmd.ExploitCmd)

	kubeConfigDir := fmt.Sprintf("%s/.kube/config", homedir.HomeDir())
	rootCmd.PersistentFlags().String("kubeconfig", kubeConfigDir, "Path to the kubeconfig file")

	rootCmd.PersistentFlags().BoolP("color", "c", false, "colored output")
	rootCmd.PersistentFlags().BoolP("insecure", "k", false, "Skip TLS verification")
	rootCmd.PersistentFlags().BoolP("no-rate-limit", "L", false, "disable rate limiting")
	rootCmd.PersistentFlags().String("as", "", "Username to impersonate")
	rootCmd.PersistentFlags().String("token", "", "Bearer token for API authentication")
	rootCmd.PersistentFlags().StringP("namespace", "n", "default", "Namespace scope")
	rootCmd.PersistentFlags().StringP("server", "s", "", "Server URL")
	rootCmd.PersistentFlags().StringP("user-agent", "A", "KAET", "HTTP User-Agent header")
}
