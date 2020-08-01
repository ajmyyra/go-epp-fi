package epp

import (
	"encoding/xml"
	"errors"
	"regexp"
	"strings"
	"time"
)

const transferKeyLowerCaseLetters = "abcdefghijklmnopqrstuvwxyz"
const transferKeyUpperCaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const transferKeyNumbers = "0123456789"
const transferKeySpecialLetters = "!\"#$%'()*+,-./:;=@[\\]^_'{|}~)"


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
		Extension *DomainExtension `xml:"extension,omitempty"`
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
	ClID      string `xml:"clID"`
	CrID      string `xml:"crID"`
	RawCrDate string `xml:"crDate"`
	CrDate    time.Time
	RawUpDate string `xml:"upDate"`
	UpDate    time.Time
	RawExDate string `xml:"exDate"`
	ExDate    time.Time
	RawTrDate string `xml:"trDate"`
	TrDate    time.Time
	AuthInfo  DomainAuthInfoResp `xml:"authInfo"`
	DsData    []DomainDSDataResp `xml:"dsData"`
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

func (s *DomainDetails) Validate() error {
	var validDomain = regexp.MustCompile(`^[a-z0-9\-]+\.fi$`)

	if !validDomain.MatchString(s.Name) {
		return errors.New("invalid domain name: " + s.Name)
	}

	if len(s.Name) < 5 || len(s.Name) > 66 {
		return errors.New("domain name must have at least 2 or a maximum of 63 characters before .fi")
	}

	if s.Registrant == "" {
		return errors.New("registrant must be defined for a domain")
	}

	if s.Period.Unit != "y" {
		return errors.New("only `y` is supported as a unit")
	}

	if s.Period.Amount < 1 || s.Period.Amount > 5 {
		return errors.New("only amount between 1-5 are supported")
	}

	return nil
}

type DomainUpdate struct {
	Xmlns string `xml:"xmlns:domain,attr"`
	Name  string `xml:"domain:name"`
	Add   struct {
		Status     *DomainStatus `xml:"domain:status,omitempty"`
		Ns         *DomainNameservers `xml:"domain:ns,omitempty"`
	} `xml:"domain:add"`
	Rem struct {
		Status     *DomainStatus       `xml:"domain:status,omitempty"`
		Ns         *DomainNameservers  `xml:"domain:ns,omitempty"`
		AuthInfo   *DomainAuthInfo `xml:"domain:authInfo,omitempty"`
	} `xml:"domain:rem"`
	Chg struct {
		Registrant string                `xml:"domain:registrant,omitempty"`
		Contact    []DomainContact       `xml:"domain:contact,omitempty"`
		AuthInfo   *DomainAuthInfo   `xml:"domain:authInfo,omitempty"`
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

type DomainAuthInfoResp struct {
	BrokerChangeKey    string `xml:"pw,omitempty"`
	OwnershipChangeKey string `xml:"pwregistranttransfer,omitempty"`
}

type DomainAuthInfo struct {
	BrokerChangeKey    string `xml:"domain:pw,omitempty"`
	OwnershipChangeKey string `xml:"domain:pwregistranttransfer,omitempty"`
}

type DomainDSDataResp struct {
	KeyTag     int `xml:"keyTag"`
	Alg        int `xml:"alg"`
	DigestType int `xml:"digestType"`
	Digest     string `xml:"digest"`
	KeyData    struct {
		Flags    int `xml:"flags"`
		Protocol int `xml:"protocol"`
		Alg      int `xml:"alg"`
		PubKey   string `xml:"pubKey"`
	} `xml:"keyData"`
}

type DomainExtension struct {
	SecDNSUpdate DomainSecDNSUpdate `xml:"secDNS:update"`
}

type DomainSecDNSUpdate struct {
	Xmlns string `xml:"xmlns:secDNS,attr"`
	Rem struct {
		DsData []DomainDSData `secDNS:dsData`
		RemoveAll bool `secDNS:all,omitempty`
	} `xml:"secDNS:rem"`
	Add struct {
		DsData []DomainDSData `secDNS:dsData`
	} `xml:"secDNS:add"`
	Chg struct {
	} `xml:"secDNS:chg"`
}

type DomainDSData struct {
	KeyTag     int `xml:"secDNS:keyTag"`
	Alg        int `xml:"secDNS:alg"`
	DigestType int `xml:"secDNS:digestType"`
	Digest     string `xml:"secDNS:digest"`
	KeyData    DomainDSKeyData `xml:"secDNS:keyData"`
}

type DomainDSKeyData struct {
	Flags    int `xml:"secDNS:flags"`
	Protocol int `xml:"secDNS:protocol"`
	Alg      int `xml:"secDNS:alg"`
	PubKey   string `xml:"secDNS:pubKey"`
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

func NewDomainUpdateSetTransferKey(domain, newKey string) (DomainUpdate, error) {
	if len(newKey) < 8 || len(newKey) > 64 {
		return DomainUpdate{}, errors.New("transfer key must be 8-64 characters long")
	}
	lowerCaseLetters := false
	upperCaseLetters := false
	specialCharacters := false
	numbers := false
	for _, c := range newKey {
		char := string(c)
		if strings.Contains(transferKeyLowerCaseLetters, char) {
			lowerCaseLetters = true
		}
		if strings.Contains(transferKeyUpperCaseLetters, char) {
			upperCaseLetters = true
		}
		if strings.Contains(transferKeySpecialLetters, char) {
			specialCharacters = true
		}
		if strings.Contains(transferKeyNumbers, char) {
			numbers = true
		}
	}

	if !lowerCaseLetters {
		return DomainUpdate{}, errors.New("transfer key does not contain a lower case letter")
	}
	if !upperCaseLetters {
		return DomainUpdate{}, errors.New("transfer key does not contain an upper case letter")
	}
	if !specialCharacters {
		return DomainUpdate{}, errors.New("transfer key does not contain a special character")
	}
	if !numbers {
		return DomainUpdate{}, errors.New("transfer key does not contain a number")
	}


	transferKeyData := createDomainUpdateBase(domain)
	transferKeyData.Chg.AuthInfo = &DomainAuthInfo{
		BrokerChangeKey: newKey,
	}

	return transferKeyData, nil
}

func NewDomainUpdateRemoveTransferKey(domain, currentKey string) DomainUpdate {
	transferKeyData := createDomainUpdateBase(domain)
	transferKeyData.Rem.AuthInfo = &DomainAuthInfo{
		BrokerChangeKey: currentKey,
	}

	return transferKeyData
}

func NewDomainSecDNSUpdate(newRecords, recordsToRemove []DomainDSData, removeAll bool) DomainSecDNSUpdate {
	secDNSUpdate := DomainSecDNSUpdate{
		Xmlns: SecDNSNamespace,
	}

	if newRecords != nil {
		secDNSUpdate.Add.DsData = newRecords
	}
	if recordsToRemove != nil {
		secDNSUpdate.Rem.DsData = recordsToRemove
	}
	secDNSUpdate.Rem.RemoveAll = removeAll

	return secDNSUpdate
}

func NewDomainDNSSecRecord(keyTag, alg, digestType int, digest string, flags, protocol, keyAlg int, pubKey string) (DomainDSData, error) {
	return DomainDSData{
		KeyTag:     keyTag,
		Alg:        alg,
		DigestType: digestType,
		Digest:     digest,
		KeyData: DomainDSKeyData{
			Flags:    flags,
			Protocol: protocol,
			Alg:      keyAlg,
			PubKey:   pubKey,
		},
	}, nil
}

func NewDomainUpdateActivateRegistryLock(domain string, numberToSend int, phoneNumbers ...string) (DomainUpdate, error)  {
	if len(phoneNumbers) < 2 || len(phoneNumbers) > 3 {
		return DomainUpdate{}, errors.New("registry lock requires 2-3 sms numbers for activation")
	}

	activationData := createDomainUpdateBase(domain)
	activationData.Chg.RegistryLock = &DomainRegistryLock{
		Type:         "activate",
		SmsNumber:    phoneNumbers,
		NumberToSend: numberToSend,
		AuthKey:      "domainauthkey",
	}

	return activationData, nil
}

func NewDomainUpdateDeactivateRegistryLock(domain, authKey string, numberToSend int, phoneNumbers ...string) DomainUpdate {
	deactivationData := createDomainUpdateBase(domain)
	deactivationData.Chg.RegistryLock = &DomainRegistryLock{
		Type:         "deactivate",
		SmsNumber:    phoneNumbers,
		NumberToSend: numberToSend,
		AuthKey:      authKey,
	}

	return deactivationData
}

func NewDomainUpdateRequestKeyForRegistryLock(domain string, numberToSend int) DomainUpdate {
	requestKeyData := createDomainUpdateBase(domain)
	requestKeyData.Chg.RegistryLock = &DomainRegistryLock{
		Type:         "requestkey",
		NumberToSend: numberToSend,
		AuthKey:      "domainauthkey",
	}

	return requestKeyData
}

func createDomainUpdateBase(domain string) DomainUpdate {
	updateData := DomainUpdate{}
	updateData.Xmlns = DomainNamespace
	updateData.Name = domain

	return updateData
}
