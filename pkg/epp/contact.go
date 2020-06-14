package epp

import (
	"encoding/xml"
)

type APIContactCheck struct {
	XMLName xml.Name `xml:"epp"`
	Xmlns   string   `xml:"xmlns,attr"`
	Command struct {
		Check struct {
			ContactCheck struct {
				Xmlns string  `xml:"xmlns:contact,attr"`
				ID      []string `xml:"contact:id"`
			} `xml:"contact:check"`
		} `xml:"check"`
		ClTRID string `xml:"clTRID"`
	} `xml:"command"`
}


type APIContactInfo struct {
	XMLName xml.Name `xml:"epp"`
	Xmlns   string   `xml:"xmlns,attr"`
	Command struct {
		Info struct {
			ContactInfo struct {
				Xmlns string `xml:"xmlns:contact,attr"`
				ID    string `xml:"contact:id"`
			} `xml:"contact:info"`
		} `xml:"info"`
		ClTRID string `xml:"clTRID"`
	} `xml:"command"`
}

type APIContactInfoResponse struct {
	XMLName  xml.Name `xml:"epp"`
	Xmlns    string   `xml:"xmlns,attr"`
	Response struct {
		Result Result `xml:"result"`
		ResData struct {
			ContactInfo ContactResponse `xml:"infData"`
		} `xml:"resData"`
		TrID Transaction `xml:"trID"`
	} `xml:"response"`
}

type APIContactCreation struct {
	XMLName xml.Name `xml:"epp"`
	Xmlns   string   `xml:"xmlns,attr"`
	Command struct {
		Create struct {
			CreateContact ContactInfo `xml:"contact:create"`
		} `xml:"create"`
		ClTRID string `xml:"clTRID"`
	} `xml:"command"`
}

type APIContactUpdate struct {
	XMLName xml.Name `xml:"epp"`
	Xmlns   string   `xml:"xmlns,attr"`
	Command struct {
		Update struct {
			ContactUpdate struct {
				Xmlns string      `xml:"xmlns:contact,attr"`
				ID    string      `xml:"contact:id"`
				Add   string      `xml:"contact:add"`
				Rem   string      `xml:"contact:rem"`
				Chg   ContactInfo `xml:"contact:chg"`
			} `xml:"contact:update"`
		} `xml:"update"`
		ClTRID string `xml:"clTRID"`
	} `xml:"command"`
}

type APIContactDeletion struct {
	XMLName xml.Name `xml:"epp"`
	Xmlns   string   `xml:"xmlns,attr"`
	Command struct {
		Delete struct {
			ContactDelete struct {
				Xmlns string `xml:"xmlns:contact,attr"`
				ID      string `xml:"contact:id"`
			} `xml:"contact:delete"`
		} `xml:"delete"`
		ClTRID string `xml:"clTRID"`
	} `xml:"command"`
}

type ContactResponse struct {
	Id         string                    `xml:"id"`
	Role       int                       `xml:"role"`
	Type       int                       `xml:"type"`
	PostalInfo ContactResponsePostalInfo `xml:"postalInfo"`
	Phone      string                    `xml:"voice"`
	Email      string                    `xml:"email"`
	LegalEmail string                    `xml:"legalemail"`
	ClID       string                    `xml:"clID"`
	CrID       string                    `xml:"crID"`
	CrDate     string                    `xml:"crDate"`
	UpDate     string                    `xml:"upDate"`
	Disclose   struct {
		DisclosedData ResponseInfoDisclosure `xml:"infDataDisclose"`
	} `xml:"disclose"`
}

type ContactResponsePostalInfo struct {
	Type           string `xml:"type,attr"`
	IsFinnish      int `xml:"isFinnish"`
	FirstName      string `xml:"firstname"`
	LastName       string `xml:"lastname"`
	Name           string `xml:"name"`
	Org            string `xml:"org"`
	BirthDate      string `xml:"birthDate"`
	Identity       string `xml:"identity"`
	RegisterNumber string `xml:"registernumber"`
	Addr           struct {
		Street     []string `xml:"street"`
		City       string   `xml:"city"`
		State      string   `xml:"sp"`
		PostalCode string   `xml:"pc"`
		Country    string   `xml:"cc"`
	} `xml:"addr"`
}

type ResponseInfoDisclosure struct {
	Flag    string `xml:"flag,attr"`
	Email   string `xml:"email"`
	Address string `xml:"address"`
}

type ContactInfo struct {
	Xmlns      string            `xml:"xmlns:contact,attr"`
	Role       int               `xml:"contact:role"`
	Type       int               `xml:"contact:type"`
	PostalInfo ContactPostalInfo `xml:"contact:postalInfo"`
	Phone      string            `xml:"contact:voice"`
	Email      string            `xml:"contact:email"`
	LegalEmail string            `xml:"contact:legalemail"`
	Disclose   ContactDisclosure `xml:"contact:disclose"`
}

type ContactPostalInfo struct {
	Type           string `xml:"type,attr"`
	IsFinnish      int `xml:"contact:isfinnish"`
	FirstName      string `xml:"contact:firstname,omitempty"`
	LastName       string `xml:"contact:lastname,omitempty"`
	Name           string `xml:"contact:name,omitempty"`
	Org            string `xml:"contact:org,omitempty"`
	BirthDate      string `xml:"contact:birthDate,omitempty"`
	Identity       string `xml:"contact:identity,omitempty"`
	RegisterNumber string `xml:"contact:registernumber,omitempty"`
	Addr           ContactUpdateAddress `xml:"contact:addr"`
}

type ContactUpdateAddress struct {
	Street     []string `xml:"contact:street"`
	City       string   `xml:"contact:city"`
	State      string   `xml:"contact:sp,omitempty"`
	PostalCode string   `xml:"contact:pc"`
	Country    string   `xml:"contact:cc"`
}

type ContactDisclosure struct {
	Flag    string `xml:"flag,attr"`
	Email   string `xml:"contact:email"`
	Address string `xml:"contact:address"`
}

func NewPrivatePersonContact(role int, finnish bool, firstName, lastName, idNumber, city, countryCode string, street []string, postalCode string, email, phone string, birthDate string) (ContactInfo, error) {
	isFinnish := 0
	if finnish {
		isFinnish = 1
	}

	contact := ContactInfo{
		Xmlns:      ContactNamespace,
		Role:       role,
		Type:       0,
		PostalInfo: ContactPostalInfo{
			Type:           "loc",
			IsFinnish:      isFinnish,
			FirstName:      firstName,
			LastName:       lastName,
			BirthDate:      birthDate,
			Identity:       idNumber,
			Addr: ContactUpdateAddress{
				Street:     street,
				City:       city,
				PostalCode: postalCode,
				Country:    countryCode,
			},
		},
		Phone:      phone,
		Email:      email,
		LegalEmail: email,
		Disclose:   ContactDisclosure{
			Flag:    "0",
			Email:   "0",
			Address: "0",
		},
	}

	return contact, nil
}

func NewBusinessContact(role int, finnish bool, orgName, registerNumber, contactName, city, countryCode string, street []string, postalCode string, email, phone string) (ContactInfo, error) {
	isFinnish := 0
	if finnish {
		isFinnish = 1
	}

	contact := ContactInfo{
		Xmlns:      ContactNamespace,
		Role:       role,
		Type:       1,
		PostalInfo: ContactPostalInfo{
			Type:           "loc",
			IsFinnish:      isFinnish,
			Name:           contactName,
			Org:            orgName,
			RegisterNumber: registerNumber,
			Addr:           ContactUpdateAddress{
				Street:     street,
				City:       city,
				PostalCode: postalCode,
				Country:    countryCode,
			},
		},
		Phone:      phone,
		Email:      email,
		LegalEmail: email,
		Disclose:   ContactDisclosure{
			Flag:    "0",
			Email:   "0",
			Address: "1",
		},
	}

	return contact, nil
}

func (s *ContactInfo) Validate() error {
	// TODO
	return nil
}