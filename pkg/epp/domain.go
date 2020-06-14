package epp

import "encoding/xml"

type APIDomainCheck struct {
	XMLName xml.Name `xml:"epp"`
	Xmlns   string   `xml:"xmlns,attr"`
	Command struct {
		Check struct {
			DomainCheck struct {
				Xmlns string   `xml:"xmlns:domain,attr"`
				Name   []string `xml:"domain:name"`
			} `xml:"domain:check"`
		} `xml:"check"`
		ClTRID string `xml:"clTRID"`
	} `xml:"command"`
}

type APIDomainInfo struct {
	XMLName xml.Name `xml:"epp"`
	Text    string   `xml:",chardata"`
	Xmlns   string   `xml:"xmlns,attr"`
	Command struct {
		Info struct {
			DomainInfo struct {
				Xmlns string `xml:"xmlns:domain,attr"`
				Name   struct {
					DomainName  string `xml:",chardata"`
					Hosts       string `xml:"hosts,attr"`
				} `xml:"domain:name"`
			} `xml:"domain:info"`
		} `xml:"info"`
		ClTRID string `xml:"clTRID"`
	} `xml:"command"`
}

type APIDomainInfoResponse struct {
	XMLName  xml.Name `xml:"epp"`
	Xmlns    string   `xml:"xmlns,attr"`
	Response struct {
		Result Result `xml:"result"`
		ResData struct {
			DomainInfo DomainInfoResp `xml:"infData"`
		} `xml:"resData"`
		TrID Transaction `xml:"trID"`
	} `xml:"response"`
}

type APIDomainCreation struct {
	XMLName  xml.Name `xml:"epp"`
	Xmlns    string   `xml:"xmlns,attr"`
	XmlnsXsi string   `xml:"xmlns:xsi,attr"`
	Command  struct {
		Create struct {
			DomainCreate DomainDetails `xml:"domain:create"`
		} `xml:"create"`
		ClTRID string `xml:"clTRID"`
	} `xml:"command"`
}

type APIDomainUpdate struct {
	XMLName xml.Name `xml:"epp"`
	Xmlns   string   `xml:"xmlns,attr"`
	Command struct {
		Update struct {
			DomainUpdate DomainUpdate `xml:"domain:update"`
		} `xml:"update"`
		ClTRID string `xml:"clTRID"`
	} `xml:"command"`
}

type DomainInfoResp struct {
	Xmlns         string `xml:"domain,attr"`
	Name          string `xml:"name"`
	RegistryLock  string `xml:"registrylock"`
	AutoRenew     string `xml:"autorenew"`
	AutoRenewDate string `xml:"autorenewDate"`
	DomainStatus  struct {
		Status    string `xml:"s,attr"`
	} `xml:"status"`
	Registrant string `xml:"registrant"`
	Contact    []DomainContact `xml:"contact"`
	Ns struct {
		HostObj []string `xml:"hostObj"`
	} `xml:"ns"`
	ClID     string `xml:"clID"`
	CrID     string `xml:"crID"`
	CrDate   string `xml:"crDate"`
	UpDate   string `xml:"upDate"`
	ExDate   string `xml:"exDate"`
	TrDate   string `xml:"trDate"`
}

type DomainDetails struct {
	Xmlns string `xml:"xmlns:domain,attr"`
	Name   string `xml:"domain:name"`
	Period struct {
		Amount int `xml:",chardata"`
		Unit string `xml:"unit,attr"`
	} `xml:"domain:period"`
	Ns DomainNameservers `xml:"domain:ns"`
	Registrant string `xml:"domain:registrant"`
	Contact    []struct {
		Text string `xml:",chardata"`
		Type string `xml:"type,attr"`
	} `xml:"domain:contact,omitempty"`
}

type DomainUpdate struct {
	Xmlns string `xml:"xmlns:domain,attr"`
	Name  string `xml:"domain:name"`
	Add   struct {
		Status DomainStatus `xml:"domain:status,omitempty"`
		Ns DomainNameservers `xml:"domain:ns,omitempty"`
	} `xml:"domain:add"`
	Rem struct {
		Status DomainStatus `xml:"domain:status,omitempty"`
		Ns DomainNameservers `xml:"domain:ns,omitempty"`
	} `xml:"domain:rem"`
	Chg struct {
		Registrant string `xml:"domain:registrant,omitempty"`
		Contact    []DomainContact `xml:"domain:contact,omitempty"`
		AuthInfo struct {
			BrokerChangeKey    string `xml:"domain:pw,omitempty"`
			OwnershipChangeKey string `xml:"domain:pwregistranttransfer,omitempty"`
		} `xml:"domain:authInfo,omitempty"`
		RegistryLock DomainRegistryLock `xml:"domain:registrylock,omitempty"`
	} `xml:"domain:chg"`
}

type DomainNameservers struct {
	HostObj  []string `xml:"domain:hostObj"`
	HostAttr []struct {
		HostName string `xml:"domain:hostName"`
		HostAddr []struct {
			Text string `xml:",chardata"`
			Ip   string `xml:"ip,attr"`
		} `xml:"domain:hostAddr"`
	} `xml:"domain:hostAttr,omitempty"`
}

type DomainStatus struct {
	Reason string `xml:",chardata"`
	Status string `xml:"s,attr"`
	Lang   string `xml:"lang,attr"`
}

type DomainContact struct {
	AccountId string `xml:",chardata"`
	Type      string `xml:"type,attr"`
}

type DomainRegistryLock struct {
	Type         string   `xml:"type,attr"`
	SmsNumber    []string `xml:"domain:smsnumber"`
	NumberToSend int      `xml:"domain:numbertosend"`
	AuthKey      string   `xml:"domain:authkey"`
}

func NewDomainDetails(domain string, years int, registrant string, dnsServers []string) DomainDetails {
	details := DomainDetails{}

	details.Xmlns = DomainNamespace
	details.Name = domain
	details.Period.Unit = "y"
	details.Period.Amount = years
	details.Registrant = registrant
	details.Ns.HostObj = dnsServers

	return details
}

func (s *DomainDetails) Validate() error {
	// TODO
	return nil
}

func DomainUpdateContacts(domain, newAdminContact, newTechContact string) DomainUpdate {
	contactData := createDomainUpdateBase(domain)

	if newAdminContact != "" {
		contactData.Chg.Contact = append(contactData.Chg.Contact,
			DomainContact{
				AccountId: newAdminContact,
				Type:      "admin",
			},
		)
	}

	if newTechContact != "" {
		contactData.Chg.Contact = append(contactData.Chg.Contact,
			DomainContact{
				AccountId: newTechContact,
				Type:      "tech",
			},
		)
	}

	return contactData
}

func DomainUpdateNameservers(domain string, removedNameservers, newNameservers []string) DomainUpdate {
	nsData := createDomainUpdateBase(domain)
	nsData.Rem.Ns.HostObj = removedNameservers
	nsData.Add.Ns.HostObj = newNameservers

	return nsData
}

func DomainUpdateSendOwnershipChangeKey(domain string) DomainUpdate {
	keyOrderData := createDomainUpdateBase(domain)
	keyOrderData.Chg.AuthInfo.OwnershipChangeKey = "new"

	return keyOrderData
}

func DomainUpdateChangeOwnership(domain string, newRegistrant, ownershipChangeKey string) DomainUpdate {
	ownershipChangeData := createDomainUpdateBase(domain)
	ownershipChangeData.Chg.Registrant = newRegistrant
	ownershipChangeData.Chg.AuthInfo.OwnershipChangeKey = ownershipChangeKey

	return ownershipChangeData
}

// TODO
// DomainUpdateActivateRegistryLock
// DomainUpdateDeactivateRegistryLock
// DomainUpdateRequestKeyForRegistryLock

func createDomainUpdateBase(domain string) DomainUpdate {
	updateData := DomainUpdate{}
	updateData.Xmlns = DomainNamespace
	updateData.Name = domain

	return updateData
}