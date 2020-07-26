package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

var msgCmd = &cobra.Command{
	Use:   "msg",
	Short: "Show and acknowledge service messages",
}

var pollShowMsgCmd = &cobra.Command{
	Use:   "show",
	Short: "Show first un-acknowledged message",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getRegistryClient(cmd)
		if err != nil {
			return err
		}

		if err = client.Connect(); err != nil {
			return errors.Wrap(err, "Unable to connect")
		}

		msg, err := client.Poll()
		if err != nil {
			return errors.Wrap(err, "Unable to fetch new message")
		}

		printJson, _ := cmd.Flags().GetBool("json")

		if printJson {
			jsonMsg, err := json.MarshalIndent(msg, "", "  ")
			if err != nil {
				return errors.Wrap(err, "Unable to create JSON message")
			}
			fmt.Println(string(jsonMsg))
		} else {
			w := tabwriter.NewWriter(os.Stdout, 2, 8, 2, ' ', 0)
			fmt.Fprintf(w, "%s\t%s\n", "ID:", msg.ID)
			fmt.Fprintf(w, "%s\t%d\n", "Count:", msg.Count)
			fmt.Fprintf(w, "%s\t%s\n", "Time:", msg.QDate)
			fmt.Fprintf(w, "%s\t%s\n", "Message:", msg.Msg)
			if msg.Name != "" {
				fmt.Fprintf(w, "%s\t%s\n", "Name:", msg.Name)
			}

			_ = w.Flush()
		}

		_ = client.Close()

		return nil
	},
}

var pollAckMsgCmd = &cobra.Command{
	Use:   "ack",
	Short: "Acknowledge specified message(s)",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("Must specify at least one message ID to acknowledge.")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getRegistryClient(cmd)
		if err != nil {
			return err
		}

		if err = client.Connect(); err != nil {
			return errors.Wrap(err, "Unable to connect")
		}

		for _, msgId := range args {
			remaining, err := client.PollAck(msgId)
			if err != nil {
				errMsg := fmt.Sprintf("Unable to ack message %s", msgId)
				return errors.Wrap(err, errMsg)
			}

			fmt.Printf("Successfully acknowledged message %s, %d messages remaining in queue.\n", msgId, remaining)
		}

		_ = client.Close()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(msgCmd)
	msgCmd.AddCommand(pollShowMsgCmd)
	msgCmd.AddCommand(pollAckMsgCmd)

	pollShowMsgCmd.Flags().BoolP("json", "j", false, "Show message as JSON")

}