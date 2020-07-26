package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to FI EPP system",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getRegistryClient(cmd)
		if err != nil {
			return err
		}

		if err = client.Connect(); err != nil {
			return errors.Wrap(err, "Unable to connect")
		}

		if err = client.Login(); err != nil {
			return errors.Wrap(err, "Unable to login")
		}

		fmt.Println("Successfully logged in.")

		_ = client.Close()

		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from FI EPP system (OBS: This cuts all sessions from the user)",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getRegistryClient(cmd)
		if err != nil {
			return err
		}

		if err = client.Connect(); err != nil {
			return errors.Wrap(err, "Unable to connect")
		}

		if err = client.Logout(); err != nil {
			return errors.Wrap(err, "Unable to logout")
		}

		fmt.Println("Successfully logged out.")

		_ = client.Close()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
}