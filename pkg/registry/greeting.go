package registry

import (
	"encoding/xml"
	"github.com/ajmyyra/go-epp-fi/pkg/epp"
	"github.com/pkg/errors"
	"time"
)

func (s *Client) Hello() (epp.Greeting, error) {
	if s.conn != nil {
		hello := epp.APIHello{
			XMLName: xml.Name{},
			Xmlns:   epp.EPPNamespace,
		}
		helloMsg, _ := xml.MarshalIndent(hello, "", "  ")
		apiResp, err := s.Send(helloMsg)
		if err != nil {
			s.logAPIConnectionError(err)
			return epp.Greeting{}, err
		}

		greeting, err := unmarshalGreeting(apiResp)
		if err != nil {
			return epp.Greeting{}, err
		}

		return greeting, nil
	}

	return epp.Greeting{}, errors.New("Uninitialized connection, unable to connect to server.")
}

func unmarshalGreeting(rawGreeting []byte) (epp.Greeting, error) {
	var greeting epp.APIGreeting
	if err := xml.Unmarshal(rawGreeting, &greeting); err != nil {
		return epp.Greeting{}, err
	}

	formattedDate, err := parseDate(greeting.Greeting.RawDate, time.RFC3339Nano)
	if err != nil {
		return epp.Greeting{}, errors.Wrap(err, "Invalid or non-existent date in greeting")
	}
	greeting.Greeting.SvDate = formattedDate

	return greeting.Greeting, nil
}
