package cmd

import (
	"github.com/Celtech/ACME/web/database"
	"github.com/Celtech/ACME/web/model"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [username] [password]",
	Short: "Adds a new API authorized user to the database",
	Long: `Adds a new API authorized user to the database.
This user will be able to fetch a access token from
the API and authorize to ALL api endpoints.

This command accepts a plain text password and will
hash it before writing to the database.

Username is preferred to be an email but anything may
be used. Note when using the API, this will be referred
to as an email regardless of whether you used an email
or not.`,
	Example: "baker-acme add support@example.com mySuperSecurePassword123",
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
