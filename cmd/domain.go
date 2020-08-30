package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Show, create, edit & delete domains",
}

var checkDomainsCmd = &cobra.Command{
	Use:   "check",
	Short: "Check domain availability",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("Must specify one or more domains")
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
		defer client.Close()

		w := tabwriter.NewWriter(os.Stdout, 2, 8, 2, ' ', 0)
		fmt.Fprintf(w, "%s\t%s\t%s\n", "Domain", "Available", "Reason")
		for _, arg := range args {
			check, err := client.CheckDomains(arg)
			if err != nil {
				return err
			}

			fmt.Fprintf(w, "%s\t%t\t%s\n", check[0].Name.Name, check[0].IsAvailable, check[0].Reason)
		}

		_ = w.Flush()
		return nil
	},
}

var showDomainCmd = &cobra.Command{
	Use:   "show",
	Short: "Show domain information",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Must specify a single domain")
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
		defer client.Close()

		domainInfo, err := client.GetDomain(args[0])
		if err != nil {
			return err
		}

		printJson, _ := cmd.Flags().GetBool("json")
		if printJson {
			jsonMsg, err := json.MarshalIndent(domainInfo, "", "  ")
			if err != nil {
				return errors.Wrap(err, "Unable to create JSON message")
			}
			fmt.Println(string(jsonMsg))

			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 2, 8, 2, ' ', 0)
		fmt.Fprintf(w, "%s\t%s\n", "Domain:", domainInfo.Name)
		fmt.Fprintf(w, "%s\t%s\n", "Status:", domainInfo.DomainStatus.Status)
		fmt.Fprintf(w, "%s\t%s\n", "Registrant:", domainInfo.Registrant)
		for _, contact := range domainInfo.Contact {
			fmt.Fprintf(w, "%s\t%s (%s)\n", "Contact:", contact.AccountId, contact.Type)
		}

		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "%s\t%s\n", "Created:", domainInfo.CrDate)
		fmt.Fprintf(w, "%s\t%s\n", "Expires:", domainInfo.ExDate)
		fmt.Fprintf(w, "%s\t%s\n", "Updated:", domainInfo.UpDate)
		fmt.Fprintf(w, "%s\t%s\n", "Transferred:", domainInfo.TrDate)

		fmt.Fprintf(w, "\n")
		autoRenew := false
		if domainInfo.AutoRenew == 1 {
			autoRenew = true
		}
		fmt.Fprintf(w, "%s\t%t\n", "Autorenew:", autoRenew)
		if autoRenew {
			fmt.Fprintf(w, "%s\t%s\n", "Autorenew date:", domainInfo.AutoRenewDate)
		}

		fmt.Fprintf(w, "\n")
		for _, ns := range domainInfo.Ns.HostObj {
			fmt.Fprintf(w, "%s\t%s\n", "Name server:", ns)
		}
		if domainInfo.AuthInfo.BrokerChangeKey != "" {
			fmt.Fprintf(w, "%s\t%s\n", "Broker change key:", domainInfo.AuthInfo.BrokerChangeKey)
		}
		if domainInfo.AuthInfo.OwnershipChangeKey != "" {
			fmt.Fprintf(w, "%s\t%s\n", "Ownership change key:", domainInfo.AuthInfo.OwnershipChangeKey)
		}

		for _, ds := range domainInfo.DsData {
			fmt.Fprintf(w, "\n")
			fmt.Fprintf(w, "%s\t%d\n", "Key tag:", ds.KeyTag)
			fmt.Fprintf(w, "%s\t%d\n", "Algorithm:", ds.Alg)
			fmt.Fprintf(w, "%s\t%d\n", "Digest type:", ds.DigestType)
			fmt.Fprintf(w, "%s\t%s\n", "Digest:", ds.Digest)
			fmt.Fprintf(w, "%s\t%d\n", "Flags:", ds.KeyData.Flags)
			fmt.Fprintf(w, "%s\t%d\n", "Key algorithm:", ds.KeyData.Alg)
			fmt.Fprintf(w, "%s\t%d\n", "Protocol:", ds.KeyData.Protocol)
			fmt.Fprintf(w, "%s\t%d\n", "Flags:", ds.KeyData.Flags)
			fmt.Fprintf(w, "%s\t%s\n", "Public key:", ds.KeyData.PubKey)
		}

		_ = w.Flush()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(domainCmd)

	domainCmd.AddCommand(checkDomainsCmd)
	domainCmd.AddCommand(showDomainCmd)

	showDomainCmd.Flags().BoolP("json", "j", false, "Show domain information as JSON")
}
