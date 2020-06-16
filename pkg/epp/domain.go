package epp

import (
	"encoding/xml"
	"time"
)

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

type APIDomainRenewal struct {
	XMLName xml.Name `xml:"epp"`
	Xmlns   string   `xml:"xmlns,attr"`
	Command struct {
		Renew struct {
			DomainRenew struct {
				Xmlns      string `xml:"xmlns:domain,attr"`
				Name       string `xml:"domain:name"`
				CurExpDate string `xml:"domain:curExpDate"`
				Period     struct {
					Years  int `xml:",chardata"`
					Unit   string `xml:"unit,attr"`
				} `xml:"domain:period"`
			} `xml:"domain:renew"`
		} `xml:"renew"`
		ClTRID string `xml:"clTRID"`
	} `xml:"command"`
}

type APIDomainTransfer struct {
	XMLName xml.Name `xml:"epp"`
	Xmlns   string   `xml:"xmlns,attr"`
	Command struct {
		Transfer struct {
			Text           string `xml:",chardata"`
			Op             string `xml:"op,attr"`
			DomainTransfer struct {
				Xmlns    string `xml:"xmlns:domain,attr"`
				Name     string `xml:"domain:name"`
				AuthInfo struct {
					TransferKey string `xml:"domain:pw"`
				} `xml:"domain:authInfo"`
				Ns *DomainNameservers `xml:"domain:ns,omitempty"`
			} `xml:"domain:transfer"`
		} `xml:"transfer"`
		ClTRID string `xml:"clTRID"`
		SvTRID string `xml:"svTRID"`
	} `xml:"command"`
}

type APIDomainDeletion struct {
	XMLName xml.Name `xml:"epp"`
	Xmlns   string   `xml:"xmlns,attr"`
	Command struct {
		Delete struct {
			DomainDelete struct {
				Xmlns      string `xml:"xmlns:domain,attr"`
				Name       string `xml:"domain:name"`
			} `xml:"domain:delete"`
		} `xml:"delete"`
		ClTRID string `xml:"clTRID"`
	} `xml:"command"`
}

type DomainInfoResp struct {
	Xmlns         string `xml:"domain,attr"`
	Name          string `xml:"name"`
	RegistryLock  string `xml:"registrylock"`
	AutoRenew     string `xml:"autorenew"`
	RawRenewDate  string `xml:"autorenewDate"`
	AutoRenewDate time.Time
	DomainStatus  struct {
		Status    string `xml:"s,attr"`
	} `xml:"status"`
	Registrant string `xml:"registrant"`
	Contact    []DomainContact `xml:"contact"`
	Ns struct {
		HostObj []string `xml:"hostObj"`
	} `xml:"ns"`
	ClID        string `xml:"clID"`
	CrID        string `xml:"crID"`
	RawCrDate   string `xml:"crDate"`
	CrDate      time.Time
	RawUpDate   string `xml:"upDate"`
	UpDate      time.Time
	RawExDate   string `xml:"exDate"`
	ExDate      time.Time
	RawTrDate   string `xml:"trDate"`
	TrDate      time.Time
	AuthInfo    DomainAuthInfo `xml:"authInfo"`
	DsData      []DomainDSData `xml:"dsData"`
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
		Status     *DomainStatus `xml:"domain:status,omitempty"`
		Ns         *DomainNameservers `xml:"domain:ns,omitempty"`
	} `xml:"domain:add"`
	Rem struct {
		Status     *DomainStatus `xml:"domain:status,omitempty"`
		Ns         *DomainNameservers `xml:"domain:ns,omitempty"`
		AuthInfo   *DomainAuthInfo `xml:"domain:authInfo,omitempty"`
	} `xml:"domain:rem"`
	Chg struct {
		Registrant string `xml:"domain:registrant,omitempty"`
		Contact    []DomainContact `xml:"domain:contact,omitempty"`
		AuthInfo   *DomainAuthInfo `xml:"domain:authInfo,omitempty"`
		RegistryLock *DomainRegistryLock `xml:"domain:registrylock,omitempty"`
	} `xml:"domain:chg"`
}

type DomainNameservers struct {
	HostObj  []string `xml:"domain:hostObj,omitempty"`
	HostAttr []struct {
		HostName string `xml:"domain:hostName,omitempty"`
		HostAddr []struct {
			Text string `xml:",chardata"`
			Ip   string `xml:"ip,attr,omitempty"`
		} `xml:"domain:hostAddr,omitempty"`
	} `xml:"domain:hostAttr,omitempty"`
}

type DomainStatus struct {
	Reason string `xml:",chardata"`
	Status string `xml:"s,attr"`
	Lang   string `xml:"lang,attr"`
}

type DomainContact struct {
	AccountId string `xml:",chardata"`
	Type      string `xml:"type,attr,omitempty"`
}

type DomainAuthInfo struct {
	BrokerChangeKey    string `xml:"domain:pw,omitempty"`
	OwnershipChangeKey string `xml:"domain:pwregistranttransfer,omitempty"`
}

type DomainDSData struct {
	KeyTag     string `xml:"keyTag"`
	Alg        string `xml:"alg"`
	DigestType string `xml:"digestType"`
	Digest     string `xml:"digest"`
	KeyData    struct {
		Text     string `xml:",chardata"`
		Flags    string `xml:"flags"`
		Protocol string `xml:"protocol"`
		Alg      string `xml:"alg"`
		PubKey   string `xml:"pubKey"`
	} `xml:"keyData"`
}

type DomainRegistryLock struct {
	Type         string   `xml:"type,attr"`
	SmsNumber    []string `xml:"domain:smsnumber,omitempty"`
	NumberToSend int      `xml:"domain:numbertosend,omitempty"`
	AuthKey      string   `xml:"domain:authkey,omitempty"`
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

func NewDomainUpdateContacts(domain, newAdminContact, newTechContact string) DomainUpdate {
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

func NewDomainUpdateNameservers(domain string, removedNameservers, newNameservers []string) DomainUpdate {
	nsData := createDomainUpdateBase(domain)
	if removedNameservers != nil {
		nsData.Rem.Ns = &DomainNameservers{
			HostObj:  removedNameservers,
		}
	}
	if newNameservers != nil {
		nsData.Add.Ns = &DomainNameservers{
			HostObj:  newNameservers,
		}
	}

	return nsData
}

func NewDomainUpdateSendOwnershipChangeKey(domain string) DomainUpdate {
	keyOrderData := createDomainUpdateBase(domain)
	keyOrderData.Chg.AuthInfo = &DomainAuthInfo{
		OwnershipChangeKey: "new",
	}

	return keyOrderData
}

func NewDomainUpdateChangeOwnership(domain string, newRegistrant, ownershipChangeKey string) DomainUpdate {
	ownershipChangeData := createDomainUpdateBase(domain)
	ownershipChangeData.Chg.Registrant = newRegistrant
	ownershipChangeData.Chg.AuthInfo = &DomainAuthInfo{
		OwnershipChangeKey: ownershipChangeKey,
	}

	return ownershipChangeData
}

func NewDomainUpdateTransferKey(domain, newKey string) DomainUpdate {
	transferKeyData := createDomainUpdateBase(domain)
	transferKeyData.Chg.AuthInfo = &DomainAuthInfo{
		BrokerChangeKey: newKey,
	}

	return transferKeyData
}

/*Transfer key validation:
* 8-64 characters.
* At least one small letter
* At least one capital letter
* At least one special character: "!\"#$%'()*+,-./:;=@[\\]^_'{|}~)"
* At least one number
*/

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