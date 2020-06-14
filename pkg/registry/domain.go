package registry

import (
	"encoding/xml"
	"github.com/ajmyyra/go-epp-fi/pkg/epp"
	"github.com/pkg/errors"
)

func (s *Client) CheckDomains(domains ...string) ([]epp.ItemCheck, error) {
	reqID := createRequestID(reqIDLength)

	domainCheck := epp.APIDomainCheck{}
	domainCheck.Xmlns = epp.EPPNamespace
	domainCheck.Command.Check.DomainCheck.Xmlns = epp.DomainNamespace
	domainCheck.Command.ClTRID = reqID

	domainCheck.Command.Check.DomainCheck.Name = domains

	checkData, err := xml.MarshalIndent(domainCheck, "", "  ")
	if err != nil {
		return []epp.ItemCheck{}, err
	}

	checkRawResp, err := s.Send(checkData)
	if err != nil {
		s.logAPIConnectionError(err, "requestID", reqID)
		return []epp.ItemCheck{}, err
	}

	var checkResult epp.APIResult
	if err = xml.Unmarshal(checkRawResp, &checkResult); err != nil {
		return []epp.ItemCheck{}, err
	}

	if checkResult.Response.Result.Code != 1000 {
		return []epp.ItemCheck{}, errors.New("Request failed: " + checkResult.Response.Result.Msg)
	}

	var checkItems []epp.ItemCheck
	for _, item := range checkResult.Response.ResData.ChkData.Cd {
		if item.Domain.Avail == 1 {
			item.IsAvailable = true
		}

		checkItems = append(checkItems, item)
	}

	return checkItems, nil
}

func (s *Client) CreateDomain(details epp.DomainDetails) (epp.CreateData, error) {
	reqID := createRequestID(reqIDLength)

	if err := details.Validate(); err != nil {
		return epp.CreateData{}, err
	}

	domainCreate := epp.APIDomainCreation{}
	domainCreate.Xmlns = epp.EPPNamespace
	domainCreate.XmlnsXsi = epp.DomainXsiNamespace
	domainCreate.Command.ClTRID = reqID

	domainCreate.Command.Create.DomainCreate = details

	createData, err := xml.MarshalIndent(domainCreate, "", "  ")
	if err != nil {
		return epp.CreateData{}, err
	}

	createRawResp, err := s.Send(createData)
	if err != nil {
		s.logAPIConnectionError(err, "requestID", reqID)
		return epp.CreateData{}, err
	}

	var createResult epp.APIResult
	if err = xml.Unmarshal(createRawResp, &createResult); err != nil {
		return epp.CreateData{}, err
	}

	if createResult.Response.Result.Code != 1000 {
		return epp.CreateData{}, errors.New("Request failed: " + createResult.Response.Result.Msg)
	}

	return createResult.Response.ResData.CreateData, nil
}

func (s *Client) GetDomain(domain string) (epp.DomainInfoResp, error) {
	reqID := createRequestID(reqIDLength)

	domainInfo := epp.APIDomainInfo{}
	domainInfo.Xmlns = epp.EPPNamespace
	domainInfo.Command.Info.DomainInfo.Xmlns = epp.DomainNamespace
	domainInfo.Command.ClTRID = reqID

	domainInfo.Command.Info.DomainInfo.Name.Hosts = "all"
	domainInfo.Command.Info.DomainInfo.Name.DomainName = domain

	infoData, err := xml.MarshalIndent(domainInfo, "", "  ")
	if err != nil {
		return epp.DomainInfoResp{}, err
	}

	infoRawResp, err := s.Send(infoData)
	if err != nil {
		s.logAPIConnectionError(err, "requestID", reqID)
		return epp.DomainInfoResp{}, err
	}

	var infoResp epp.APIDomainInfoResponse
	if err = xml.Unmarshal(infoRawResp, &infoResp); err != nil {
		return epp.DomainInfoResp{}, err
	}

	if infoResp.Response.Result.Code != 1000 {
		return epp.DomainInfoResp{}, errors.New("Request failed: " + infoResp.Response.Result.Msg)
	}

	return infoResp.Response.ResData.DomainInfo, nil
}

func (s *Client) UpdateDomain(update epp.DomainUpdate) error {
	reqID := createRequestID(reqIDLength)

	domainUpdate := epp.APIDomainUpdate{}
	domainUpdate.Xmlns = epp.EPPNamespace
	domainUpdate.Command.ClTRID = reqID

	domainUpdate.Command.Update.DomainUpdate = update

	updateData, err := xml.MarshalIndent(domainUpdate, "", "  ")
	if err != nil {
		return err
	}

	updateRawResp, err := s.Send(updateData)
	if err != nil {
		s.logAPIConnectionError(err, "requestID", reqID)
		return err
	}

	var updateResp epp.APIResult
	if err = xml.Unmarshal(updateRawResp, &updateResp); err != nil {
		return err
	}

	if updateResp.Response.Result.Code != 1000 {
		return errors.New("Request failed: " + updateResp.Response.Result.Msg)
	}

	return nil
}

func (s *Client) RenewDomain(domain string, years int) error {
	// TODO
	return nil
}

func (s *Client) TransferDomain(domain, transferKey string) error {
	// TODO
	return nil
}

func (s *Client) DeleteDomain(domain string) error {
	// TODO
	return nil
}