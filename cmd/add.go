package cmd

import (
	"baker-acme/web/database"
	"baker-acme/web/model"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [username] [password]",
	Short: "Adds a new API authorized user to the database",
	Long:  "Adds a new API authorized user to the database",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			log.Error("username or password argument is missing")
			return
		}
		var user = model.User{
			Email:    args[0],
			Password: args[1],
		}
		database.GetDB().Create(&user)

		log.Infof("User Added:\r\n\r\nEmail: %s\r\nPassword: %s", user.Email, user.Password)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
