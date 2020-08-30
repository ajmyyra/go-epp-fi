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

var contactCmd = &cobra.Command{
	Use:   "contact",
	Short: "Show, create, edit & delete contacts",
}

var createContactCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new contact",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getRegistryClient(cmd)
		if err != nil {
			return err
		}

		contactInfo, err := createContactInfo(cmd)
		if err != nil {
			return err
		}

		if err = client.Connect(); err != nil {
			return errors.Wrap(err, "Unable to connect")
		}
		defer client.Close()

		contactId, err := client.CreateContact(contactInfo)
		if err != nil {
			return err
		}

		fmt.Printf("Created a new contact with ID: %s\n", contactId)

		return nil
	},
}

var updateContactCmd = &cobra.Command{
	Use:   "update",
	Short: "Update contact information",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("Must specify contact id")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getRegistryClient(cmd)
		if err != nil {
			return err
		}

		contactId := args[0]

		contactInfo, err := createContactInfo(cmd)
		if err != nil {
			return err
		}

		if err = client.Connect(); err != nil {
			return errors.Wrap(err, "Unable to connect")
		}
		defer client.Close()

		if err = client.UpdateContact(contactId, contactInfo); err != nil {
			return err
		}

		fmt.Println("Contact updated successfully.")

		return nil
	},
}

var deleteContactCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete contact",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("Must specify contact id")
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

		if err = client.DeleteContact(args[0]); err != nil {
			return errors.Wrap(err, "Unable to delete contact")
		}

		return nil
	},
}

var checkContactsCmd = &cobra.Command{
	Use:   "check",
	Short: "Check contact availability",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("Must specify one or more contact id's")
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
		fmt.Fprintf(w, "%s\t%s\t%s\n", "Contact", "Available", "Reason")
		for _, arg := range args {
			check, err := client.CheckContacts(arg)
			if err != nil {
				return err
			}

			fmt.Fprintf(w, "%s\t%t\t%s\n", check[0].Id.Name, check[0].IsAvailable, check[0].Reason)
		}

		_ = w.Flush()
		return nil
	},
}

var showContactCmd = &cobra.Command{
	Use:   "show",
	Short: "Show information for a specific account",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("Must specify contact id")
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

		contact, err := client.GetContact(args[0])
		if err != nil {
			return errors.Wrap(err, "Unable to fetch contact information")
		}

		printJson, _ := cmd.Flags().GetBool("json")

		if printJson {
			jsonMsg, err := json.MarshalIndent(contact, "", "  ")
			if err != nil {
				return errors.Wrap(err, "Unable to create JSON message")
			}
			fmt.Println(string(jsonMsg))
		} else {
			w := tabwriter.NewWriter(os.Stdout, 2, 8, 2, ' ', 0)
			fmt.Fprintf(w, "%s\t%s\n", "ID:", contact.Id)
			fmt.Fprintf(w, "%s\t%d\n", "Type:", contact.Type)
			fmt.Fprintf(w, "%s\t%d\n", "Role:", contact.Role)
			finnish := "yes"
			if contact.PostalInfo.IsFinnish != 1 {
				finnish = "no"
			}
			fmt.Fprintf(w, "%s\t%s\n", "Finnish:", finnish)
			if contact.Type == 0 {
				name := fmt.Sprintf("%s %s", contact.PostalInfo.FirstName, contact.PostalInfo.LastName)
				fmt.Fprintf(w, "%s\t%s\n", "Name:", name)
				if contact.PostalInfo.Identity != "" {
					fmt.Fprintf(w, "%s\t%s\n", "Identity:", contact.PostalInfo.Identity)
				}
			} else {
				fmt.Fprintf(w, "%s\t%s\n", "Name:", contact.PostalInfo.Name)
				fmt.Fprintf(w, "%s\t%s\n", "Org:", contact.PostalInfo.Org)
				fmt.Fprintf(w, "%s\t%s\n", "VAT number:", contact.PostalInfo.RegisterNumber)
			}

			fmt.Fprintf(w, "\n")
			for _, street := range contact.PostalInfo.Addr.Street {
				fmt.Fprintf(w, "%s\t%s\n", "Street:", street)
			}
			fmt.Fprintf(w, "%s\t%s\n", "Postal code:", contact.PostalInfo.Addr.PostalCode)
			fmt.Fprintf(w, "%s\t%s\n", "City:", contact.PostalInfo.Addr.City)
			if contact.PostalInfo.Addr.State != "" {
				fmt.Fprintf(w, "%s\t%s\n", "State:", contact.PostalInfo.Addr.State)
			}
			fmt.Fprintf(w, "%s\t%s\n", "Country:", contact.PostalInfo.Addr.Country)

			fmt.Fprintf(w, "\n")
			email := contact.Email
			if email == "" {
				email = contact.LegalEmail
			}
			fmt.Fprintf(w, "%s\t%s\n", "Email:", email)
			fmt.Fprintf(w, "%s\t%s\n", "Phone:", contact.Phone)

			fmt.Fprintf(w, "\n")
			fmt.Fprintf(w, "%s\t%s\n", "Created:", contact.CrDate)
			fmt.Fprintf(w, "%s\t%s\n", "Updated:", contact.UpDate)

			_ = w.Flush()
		}

		return nil
	},
}

func createContactInfo(cmd *cobra.Command) (epp.ContactInfo, error) {
	contactType, _ := cmd.Flags().GetInt("type")

	contactRole, _ := cmd.Flags().GetInt("role")
	finnish := true
	nonFinnish, _ := cmd.Flags().GetBool("non-finnish")
	if nonFinnish {
		finnish = false
	}
	firstName, _ := cmd.Flags().GetString("firstname")
	lastName, _ := cmd.Flags().GetString("lastname")
	identity, _ := cmd.Flags().GetString("identity")
	birthDate, _ := cmd.Flags().GetString("birthdate")

	registerNumber, _ := cmd.Flags().GetString("registernumber")
	orgName, _ := cmd.Flags().GetString("org")

	streetAddr, _ := cmd.Flags().GetString("streetaddress")
	postalCode, _ := cmd.Flags().GetString("postalcode")
	city, _ := cmd.Flags().GetString("city")
	country, _ := cmd.Flags().GetString("country")
	email, _ := cmd.Flags().GetString("email")
	phone, _ := cmd.Flags().GetString("phone")

	var contactInfo epp.ContactInfo
	var err error
	if contactType == 0 {
		contactInfo, err = epp.NewPrivatePersonContact(
			contactRole,
			finnish,
			firstName,
			lastName,
			identity,
			city,
			country,
			[]string{streetAddr},
			postalCode,
			email,
			phone,
			birthDate,
		)
	} else {
		contactInfo, err = epp.NewBusinessContact(
			contactRole,
			finnish,
			orgName,
			registerNumber,
			fmt.Sprintf("%s %s", firstName, lastName),
			city,
			country,
			[]string{streetAddr},
			postalCode,
			email,
			phone,
		)

		contactInfo.Type = contactType
	}

	if err != nil {
		return epp.ContactInfo{}, err
	}

	err = contactInfo.Validate()
	if err != nil {
		return epp.ContactInfo{}, err
	}

	return contactInfo, nil
}

func populateContactManagementFlags(c *cobra.Command) {
	c.Flags().Int("type", -1, "Contact type")
	c.Flags().Int("role", -1, "Contact role")
	c.Flags().Bool("non-finnish", false, "Contact is not from Finland")

	c.Flags().String("firstname", "", "First name")
	c.Flags().String("lastname", "", "Last name")
	c.Flags().String("identity", "", "Identity")
	c.Flags().String("birthdate", "", "Birth date (YYYY-MM-DD)")

	c.Flags().String("org", "", "Organisation name")
	c.Flags().String("registernumber", "", "Organisation/company register number")

	c.Flags().String("streetaddress", "", "Street address")
	c.Flags().String("postalcode", "", "Postal code")
	c.Flags().String("city", "", "City")
	c.Flags().String("country", "", "Country code (FI, US etc)")

	c.Flags().String("email", "", "Email")
	c.Flags().String("phone", "", "Phone")
}

func init() {
	rootCmd.AddCommand(contactCmd)
	contactCmd.AddCommand(checkContactsCmd)
	contactCmd.AddCommand(showContactCmd)
	contactCmd.AddCommand(createContactCmd)
	contactCmd.AddCommand(updateContactCmd)
	contactCmd.AddCommand(deleteContactCmd)

	showContactCmd.Flags().BoolP("json", "j", false, "Show contact information as JSON")

	populateContactManagementFlags(createContactCmd)
	populateContactManagementFlags(updateContactCmd)
}
