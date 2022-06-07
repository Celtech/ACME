package cmd

import (
	"baker-acme/web/model"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// hashCmd represents the hash command
var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: "Returns a hashed version of a plaintext password",
	Long:  "Returns a hashed version of a plaintext password",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Error("password argument is missing")
			return
		}

		password, err := model.Hash(args[0])
		if err != nil {
			log.Error(err)
			return
		}

		log.Infof("Hashed string: %s", string(password))
	},
}

func init() {
	rootCmd.AddCommand(hashCmd)
}
