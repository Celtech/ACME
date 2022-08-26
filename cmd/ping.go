package cmd

import (
	"fmt"
	acmeConfig "github.com/Celtech/ACME/config"
	"github.com/spf13/cobra"
	"net/http"
	"os"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Returns the health status of your server",
	Long: `Makes a health check request to query the status
of your web server. It will return "ALIVE" with
the exit code 0 when healthy or "DEAD" with exit
code 1 when unhealthy.`,
	Run: func(cmd *cobra.Command, args []string) {
		serverPort := acmeConfig.GetConfig().GetString("server.port")

		resp, err := http.Get("http://127.0.0.1:" + serverPort + "/ping") // Note pointer dereference using *

		// If there is an error or non-200 status, exit with 1 signaling unsuccessful check.
		if err != nil || resp.StatusCode != 200 {
			fmt.Println("DEAD")
			os.Exit(1)
		}
		fmt.Println("ALIVE")
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
