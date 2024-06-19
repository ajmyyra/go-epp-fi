package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/ajmyyra/go-epp-fi/pkg/epp"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Show, create, edit, transfer & delete domains",
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

var registerDomainCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new domain",
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

		years, _ := cmd.Flags().GetInt("years")
		registrant, _ := cmd.Flags().GetString("registrant")
		ns, _ := cmd.Flags().GetStringArray("nameserver")

		domainDetails, err := epp.NewDomainDetails(args[0], years, registrant, ns)
		if err != nil {
			return err
		}
		if err = domainDetails.Validate(); err != nil {
			return err
		}

		if err = client.Connect(); err != nil {
			return errors.Wrap(err, "Unable to connect")
		}
		defer client.Close()

		createdDomain, err := client.CreateDomain(domainDetails)
		if err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 2, 8, 2, ' ', 0)
		fmt.Fprintf(w, "%s\t%s\n", "Registered domain:", createdDomain.Name)
		fmt.Fprintf(w, "%s\t%s\n", "Created:", createdDomain.CrDate)
		fmt.Fprintf(w, "%s\t%s\n", "Expires at:", createdDomain.ExDate)
		_ = w.Flush()

		return nil
	},
}

var renewDomainCmd = &cobra.Command{
	Use:   "renew",
	Short: "Renews an existing domain",
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

		years, _ := cmd.Flags().GetInt("years")
		currentExp, _ := cmd.Flags().GetString("expiration")
		if currentExp == "" {
			return errors.New("Current expiration date must be defined.")
		}

		if err = client.Connect(); err != nil {
			return errors.Wrap(err, "Unable to connect")
		}
		defer client.Close()

		renewedDomain, err := client.RenewDomain(args[0], currentExp, years)
		if err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 2, 8, 2, ' ', 0)
		fmt.Fprintf(w, "%s\t%s\n", "Renewed domain:", renewedDomain.Name)
		fmt.Fprintf(w, "%s\t%s\n", "New expiration date:", renewedDomain.ExDate)
		_ = w.Flush()

		return nil
	},
}

var updateDomainNameserversCmd = &cobra.Command{
	Use:   "update-ns",
	Short: "Update name servers for the domain",
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

		removedNs, _ := cmd.Flags().GetStringArray("remove-ns")
		addNs, _ := cmd.Flags().GetStringArray("add-ns")
		if len(removedNs) == 0 && len(addNs) == 0 {
			return errors.New("Must add or remove one or more nameservers.")
		}

		if err = client.Connect(); err != nil {
			return errors.Wrap(err, "Unable to connect")
		}
		defer client.Close()

		domainUpdate := epp.NewDomainUpdateNameservers(args[0], removedNs, addNs)
		if err = client.UpdateDomain(domainUpdate); err != nil {
			return err
		}

		fmt.Printf("Name servers for domain %s updated successfully.\n", args[0])

		return nil
	},
}

var transferDomainCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer domains, set & remove transfer keys",
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

		key, _ := cmd.Flags().GetString("key")
		if key == "" {
			return errors.New("Transfer key must be specified.")
		}
		newNs, _ := cmd.Flags().GetStringArray("new-ns")

		if err = client.Connect(); err != nil {
			return errors.Wrap(err, "Unable to connect")
		}
		defer client.Close()

		transferData, err := client.TransferDomain(args[0], key, newNs)
		if err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 2, 8, 2, ' ', 0)
		fmt.Fprintf(w, "%s\t%s\n", "Transferred domain:", transferData.Name)
		fmt.Fprintf(w, "%s\t%s\n", "Transfer status:", transferData.TrStatus)
		fmt.Fprintf(w, "%s\t%s\n", "Transfer date:", transferData.ReDate)

		return nil
	},
}

var setDomainTransferKeyCmd = &cobra.Command{
	Use:   "set-key",
	Short: "Set transfer key for domain",
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

		key, _ := cmd.Flags().GetString("key")
		if key == "" {
			return errors.New("Transfer key must be specified.")
		}

		transferKeyUpdate, err := epp.NewDomainUpdateSetTransferKey(args[0], key)
		if err != nil {
			return err
		}

		if err = client.Connect(); err != nil {
			return errors.Wrap(err, "Unable to connect")
		}
		defer client.Close()

		if err = client.UpdateDomain(transferKeyUpdate); err != nil {
			return err
		}

		return nil
	},
}

var removeDomainTransferKeyCmd = &cobra.Command{
	Use:   "remove-key",
	Short: "Remove transfer key for domain",
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

		transferKeyRemoval := epp.NewDomainUpdateRemoveTransferKey(args[0])

		if err = client.Connect(); err != nil {
			return errors.Wrap(err, "Unable to connect")
		}
		defer client.Close()

		if err = client.UpdateDomain(transferKeyRemoval); err != nil {
			return err
		}

		return nil
	},
}

var deleteDomain = &cobra.Command{
	Use:   "delete",
	Short: "Delete existing domain",
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

		if err = client.DeleteDomain(args[0]); err != nil {
			return err
		}

		fmt.Printf("Domain %s successfully deleted.\n", args[0])

		return nil
	},
}

func init() {
	rootCmd.AddCommand(domainCmd)

	domainCmd.AddCommand(checkDomainsCmd)
	domainCmd.AddCommand(showDomainCmd)
	domainCmd.AddCommand(registerDomainCmd)
	domainCmd.AddCommand(renewDomainCmd)
	domainCmd.AddCommand(updateDomainNameserversCmd)
	domainCmd.AddCommand(deleteDomain)

	domainCmd.AddCommand(transferDomainCmd)
	transferDomainCmd.AddCommand(setDomainTransferKeyCmd)
	transferDomainCmd.AddCommand(removeDomainTransferKeyCmd)

	showDomainCmd.Flags().BoolP("json", "j", false, "Show domain information as JSON")

	registerDomainCmd.Flags().IntP("years", "y", 1, "Domain registration period (1-5)")
	registerDomainCmd.Flags().String("registrant", "", "Domain registrant ID")
	registerDomainCmd.Flags().StringArray("ns", []string{}, "Domain name servers")

	renewDomainCmd.Flags().IntP("years", "y", 1, "Domain renewal period (1-5)")
	renewDomainCmd.Flags().String("expiration", "", "Current expiration date (YYYY-MM-DD)")

	updateDomainNameserversCmd.Flags().StringArray("remove-ns", []string{}, "Name server to remove (can be specified more than once)")
	updateDomainNameserversCmd.Flags().StringArray("add-ns", []string{}, "Name server to add (can be specified more than once)")

	transferDomainCmd.Flags().String("key", "", "Domain transfer key")
	transferDomainCmd.Flags().StringArray("new-ns", []string{}, "Set new name servers for the transferred domain")
	setDomainTransferKeyCmd.Flags().String("key", "", "New domain transfer key")
}
