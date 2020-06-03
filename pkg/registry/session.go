package registry

import (
	"encoding/xml"
	"github.com/ajmyyra/go-epp-fi/pkg/epp"
	"github.com/pkg/errors"
)

func (s *Client) Login() error {
	loginDetails := epp.Login{}
	loginDetails.ClID = s.Credentials.Username
	loginDetails.Pw = s.Credentials.Password

	loginDetails.Options.Version = APIVersion
	loginDetails.Options.Lang = APILanguage

	loginDetails.Svcs.ObjURI = []string{epp.DomainNamespace, epp.HostNamespace, epp.ContactNamespace}
	loginDetails.Svcs.SvcExtension.ExtURI = []string{epp.SecDNSNamespace, epp.DomainExtNamespace}

	EPPLogin := epp.APILogin{}
	EPPLogin.Xmlns = epp.EPPNamespace
	EPPLogin.Command.Login = loginDetails
	EPPLogin.Command.ClTRID = createRequestID(reqIDLength)

	loginData, err := xml.MarshalIndent(EPPLogin, "", "  ")
	if err != nil {
		return errors.Wrap(err, "Problem converting login message to XML")
	}

	rawResult, err := s.Send(loginData)
	if err != nil {
		return errors.Wrap(err, "Login failed")
	}

	loginResult := epp.APIResult{}
	if err = xml.Unmarshal(rawResult, &loginResult); err != nil {
		return errors.Wrap(err, "Unrecognised result body")
	}

	result := loginResult.Response.Result
	if result.Code != 1000 {
		return errors.New(result.Msg)
	}

	s.LoggedIn = true
	return nil
}

func (s *Client) Logout() error {
	EPPLogout := epp.APILogout{}
	EPPLogout.Xmlns = epp.EPPNamespace
	EPPLogout.Command.ClTRID = createRequestID(reqIDLength)

	logoutData, err := xml.MarshalIndent(EPPLogout, "", "  ")
	if err != nil {
		return errors.Wrap(err, "Problem converting logout message to XML")
	}

	rawResult, err := s.Send(logoutData)
	if err != nil {
		return errors.Wrap(err, "Logout failed")
	}

	logoutResult := epp.APIResult{}
	if err = xml.Unmarshal(rawResult, &logoutResult); err != nil {
		return errors.Wrap(err,"Unrecognised result body")
	}

	result := logoutResult.Response.Result
	if result.Code != 1500 {
		return errors.New(result.Msg)
	}

	s.LoggedIn = false
	return nil
}

func (s *Client) ChangePassword() error {
	// TODO
	return nil
}