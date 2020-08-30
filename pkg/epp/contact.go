package epp

import (
	"encoding/xml"
	"github.com/pkg/errors"
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
	Id         string                    `xml:"id" json:"id"`
	Role       int                       `xml:"role" json:"role"`
	Type       int                       `xml:"type" json:"type"`
	PostalInfo ContactResponsePostalInfo `xml:"postalInfo" json:"postal_info"`
	Phone      string                    `xml:"voice" json:"voice"`
	Email      string                    `xml:"email" json:"email"`
	LegalEmail string                    `xml:"legalemail" json:"legal_email"`
	ClID       string                    `xml:"clID" json:"-"`
	CrID       string                    `xml:"crID" json:"creator"`
	CrDate     string                    `xml:"crDate" json:"creation_date"`
	UpDate     string                    `xml:"upDate" json:"update_date"`
	Disclose   struct {
		DisclosedData ResponseInfoDisclosure `xml:"infDataDisclose" json:"information_disclosure"`
	} `xml:"disclose" json:"disclosure"`
}

type ContactResponsePostalInfo struct {
	Type           string `xml:"type,attr" json:"-"`
	IsFinnish      int `xml:"isFinnish" json:"is_finnish"`
	FirstName      string `xml:"firstname" json:"first_name"`
	LastName       string `xml:"lastname" json:"last_name"`
	Name           string `xml:"name" json:"name"`
	Org            string `xml:"org" json:"org"`
	BirthDate      string `xml:"birthDate" json:"birth_date"`
	Identity       string `xml:"identity" json:"identity"`
	RegisterNumber string `xml:"registernumber" json:"register_number"`
	Addr           struct {
		Street     []string `xml:"street" json:"street"`
		City       string   `xml:"city" json:"city"`
		State      string   `xml:"sp" json:"state"`
		PostalCode string   `xml:"pc" json:"postal_code"`
		Country    string   `xml:"cc" json:"country"`
	} `xml:"addr"`
}

type ResponseInfoDisclosure struct {
	Flag    string `xml:"flag,attr" json:"flag"`
	Email   string `xml:"email" json:"email"`
	Address string `xml:"address" json:"address"`
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
	PostalInfoType string               `xml:"type,attr"`
	IsFinnish      int                  `xml:"contact:isfinnish"`
	FirstName      string               `xml:"contact:firstname,omitempty"`
	LastName       string               `xml:"contact:lastname,omitempty"`
	Name           string               `xml:"contact:name,omitempty"`
	Org            string               `xml:"contact:org,omitempty"`
	BirthDate      string               `xml:"contact:birthDate,omitempty"`
	Identity       string               `xml:"contact:identity,omitempty"`
	RegisterNumber string               `xml:"contact:registernumber,omitempty"`
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
	Flag    int `xml:"flag,attr"`
	Email   int `xml:"contact:email"`
	Address int `xml:"contact:address"`
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
			PostalInfoType: "loc",
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
			Flag:    0,
			Email:   0,
			Address: 0,
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
			PostalInfoType: "loc",
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
			Flag:    0,
			Email:   0,
			Address: 1,
		},
	}

	return contact, nil
}

func (s *ContactInfo) Validate() error {
	if s.PostalInfo.IsFinnish != 0 && s.PostalInfo.IsFinnish != 1 {
		return errors.New("IsFinnish attribute is a numeric boolean, so it must be either 1 or 0.")
	}
	if s.PostalInfo.PostalInfoType != "loc" {
		return errors.New("Type attribute for postal info must be 'loc'.")
	}

	if s.Role < 2 || s.Role > 5 {
		return errors.New("Role must be between 2-5. See README.md for different roles.")
	}
	if s.Role == 5 {
		if s.LegalEmail == "" {
			return errors.New("LegalEmail must be defined for registrants (role 5).")
		}
	} else {
		if s.Email == "" {
			return errors.New("Email must be defined for non-registrants (role != 5)")
		}
	}

	if s.Type == 0 {
		if s.PostalInfo.FirstName == "" || s.PostalInfo.LastName == "" {
			return errors.New("Private person (type 0) must have FirstName and LastName.")
		}
		if len(s.PostalInfo.FirstName) > 255 {
			return errors.New("FirstName must have less than 255 characters.")
		}
		if len(s.PostalInfo.LastName) > 255 {
			return errors.New("LastName must have less than 255 characters.")
		}

		if s.PostalInfo.Name != "" {
			return errors.New("Private person's name is defined as FirstName and LastName.")
		}

		if s.PostalInfo.IsFinnish == 1 {
			if s.PostalInfo.Identity == "" {
				return errors.New("Finnish person must have their national identification number (hetu) defined as Identity.")
			}
		} else {
			if s.PostalInfo.BirthDate == "" {
				return errors.New("Non-finnish person must have their birthdate defined as YYYY-MM-DD.")
			}
		}
	} else if s.Type >= 1 && s.Type <= 7 {
		if s.PostalInfo.FirstName != "" || s.PostalInfo.LastName != "" {
			return errors.New("Organisations may specify their contacts/departments name as Name.")
		}
		if len(s.PostalInfo.Name) > 255 {
			return errors.New("Name must have less than 255 characters.")
		}

		if s.PostalInfo.Org == "" {
			return errors.New("Organisations must specify their name as Org.")
		}
		if len(s.PostalInfo.Org) > 255 {
			return errors.New("Org must be 2-255 characters long.")
		}

		if s.PostalInfo.RegisterNumber == "" {
			return errors.New("Organisation must specify a valid RegisterNumber.")
		}

		if s.Disclose.Address == 0 {
			return errors.New("Organisations cannot disclose their address.")
		}
	} else {
		return errors.New("Type must be between 0-7. See README.md for different types.")
	}

	if len(s.PostalInfo.Addr.Street) < 1 || len(s.PostalInfo.Addr.Street) > 3 {
		return errors.New("Street must have 1-3 items within the array.")
	}
	for _, street := range s.PostalInfo.Addr.Street {
		if len(street) < 2 || len(street) > 255 {
			return errors.New("Street array members must be 2-255 characters long.")
		}
	}

	if len(s.PostalInfo.Addr.City) < 2 || len(s.PostalInfo.Addr.City) > 128 {
		return errors.New("City must be 2-128 characters long.")
	}

	if len(s.PostalInfo.Addr.State) > 128 {
		return errors.New("State must be 2-128 characters long, if defined.")
	}

	if s.PostalInfo.IsFinnish == 1 {
		if len(s.PostalInfo.Addr.PostalCode) != 5 {
			return errors.New("Finnish postal code must be 5 characters long.")
		}
	} else {
		if len(s.PostalInfo.Addr.PostalCode) < 2 || len(s.PostalInfo.Addr.PostalCode) > 16 {
			return errors.New("Non-Finnish postal code must be 2-16 characters long.")
		}
	}

	if len(s.PostalInfo.Addr.Country) != 2 {
		return errors.New("Country code must be ISO 3166-1 alpha-2 -formatted (2 characters long).")
	}

	correctPhoneNumber := true
	if len(s.Phone) < 5 {
		correctPhoneNumber = false
	} else {
		if string(s.Phone[0]) != "+" {
			correctPhoneNumber = false
		}
	}

	if !correctPhoneNumber {
		return errors.New("Phone number must be defined with country code, e.g. +358401234567")
	}

	return nil
}
