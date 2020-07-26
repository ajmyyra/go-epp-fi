package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Show account balance",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getRegistryClient(cmd)
		if err != nil {
			return err
		}

		if err = client.Connect(); err != nil {
			return errors.Wrap(err, "Unable to connect")
		}

		balance, err := client.Balance()
		if err != nil {
			return errors.Wrap(err, "Unable to fetch account balance")
		}

		fmt.Printf("Account currently has %d euros available.\n", balance)

		_ = client.Close()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(balanceCmd)
}