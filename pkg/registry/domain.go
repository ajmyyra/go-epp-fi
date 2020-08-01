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
		if item.Name.Avail == 1 {
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

	createDataResp := createResult.Response.ResData.CreateData
	createDataResp.CrDate, err = parseDate(createDataResp.RawCrDate)
	if err != nil {
		return createDataResp, err
	}
	createDataResp.ExDate, err = parseDate(createDataResp.RawExDate)
	if err != nil {
		return createDataResp, err
	}

	return createDataResp, nil
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

	domInfo := infoResp.Response.ResData.DomainInfo

	domInfo.AutoRenewDate, err = parseDate(domInfo.RawRenewDate)
	if err != nil {
		return epp.DomainInfoResp{}, err
	}

	domInfo.CrDate, err = parseDate(domInfo.RawCrDate)
	if err != nil {
		return epp.DomainInfoResp{}, err
	}

	domInfo.UpDate, err = parseDate(domInfo.RawUpDate)
	if err != nil {
		return epp.DomainInfoResp{}, err
	}

	domInfo.ExDate, err = parseDate(domInfo.RawExDate)
	if err != nil {
		return epp.DomainInfoResp{}, err
	}

	domInfo.TrDate, err = parseDate(domInfo.RawTrDate)
	if err != nil {
		return epp.DomainInfoResp{}, err
	}

	return domInfo, nil
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

func (s *Client) UpdateDomainExtensions(extUpdate epp.DomainExtension) error {
	reqID := createRequestID(reqIDLength)

	domainUpdate := epp.APIDomainUpdate{}
	domainUpdate.Xmlns = epp.EPPNamespace
	domainUpdate.Command.ClTRID = reqID

	domainUpdate.Command.Extension = &extUpdate

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

func (s *Client) RenewDomain(domain, currentExpiration string, years int) (epp.RenewalData, error) {
	reqID := createRequestID(reqIDLength)

	domainRenewal := epp.APIDomainRenewal{}
	domainRenewal.Xmlns = epp.EPPNamespace
	domainRenewal.Command.Renew.DomainRenew.Xmlns = epp.DomainNamespace
	domainRenewal.Command.ClTRID = reqID

	domainRenewal.Command.Renew.DomainRenew.Name = domain
	domainRenewal.Command.Renew.DomainRenew.CurExpDate = currentExpiration
	domainRenewal.Command.Renew.DomainRenew.Period.Unit = "y"
	domainRenewal.Command.Renew.DomainRenew.Period.Years = years

	renewalData, err := xml.MarshalIndent(domainRenewal, "", "  ")
	if err != nil {
		return epp.RenewalData{}, err
	}

	renewRawResp, err := s.Send(renewalData)
	if err != nil {
		s.logAPIConnectionError(err, "requestID", reqID)
		return epp.RenewalData{}, err
	}

	var renewResp epp.APIResult
	if err = xml.Unmarshal(renewRawResp, &renewResp); err != nil {
		return epp.RenewalData{}, err
	}

	if renewResp.Response.Result.Code != 1000 {
		return epp.RenewalData{}, errors.New("Request failed: " + renewResp.Response.Result.Msg)
	}

	renewalInfo := renewResp.Response.ResData.RenewalData
	renewalInfo.ExpireDate, err = parseDate(renewalInfo.RawExpDate)
	if err != nil {
		return epp.RenewalData{}, err
	}

	return renewalInfo, nil
}

func (s *Client) TransferDomain(domain, transferKey string, newNameservers []string) (epp.TransferData, error) {
	reqID := createRequestID(reqIDLength)

	domainTransfer := epp.APIDomainTransfer{}
	domainTransfer.Xmlns = epp.EPPNamespace
	domainTransfer.Command.Transfer.DomainTransfer.Xmlns = epp.DomainNamespace
	domainTransfer.Command.ClTRID = reqID

	domainTransfer.Command.Transfer.Op = "request"
	domainTransfer.Command.Transfer.DomainTransfer.Name = domain
	domainTransfer.Command.Transfer.DomainTransfer.AuthInfo.TransferKey = transferKey

	if newNameservers != nil {
		domainTransfer.Command.Transfer.DomainTransfer.Ns = &epp.DomainNameservers{
			HostObj:  newNameservers,
		}
	}

	transferData, err := xml.MarshalIndent(domainTransfer, "", "  ")
	if err != nil {
		return epp.TransferData{}, err
	}

	transferRawResp, err := s.Send(transferData)
	if err != nil {
		s.logAPIConnectionError(err, "requestID", reqID)
		return epp.TransferData{}, err
	}

	var transferResp epp.APIResult
	if err = xml.Unmarshal(transferRawResp, &transferResp); err != nil {
		return epp.TransferData{}, err
	}

	if transferResp.Response.Result.Code != 1000 {
		return epp.TransferData{}, errors.New("Request failed: " + transferResp.Response.Result.Msg)
	}

	transfer := transferResp.Response.ResData.TransferData
	transfer.ReDate, err = parseDate(transfer.ReRawDate)
	if err != nil {
		return epp.TransferData{}, err
	}

	return transfer, nil
}

func (s *Client) DeleteDomain(domain string) error {
	reqID := createRequestID(reqIDLength)

	domainDeletion := epp.APIDomainDeletion{}
	domainDeletion.Xmlns = epp.EPPNamespace
	domainDeletion.Command.Delete.DomainDelete.Xmlns = epp.DomainNamespace
	domainDeletion.Command.ClTRID = reqID

	domainDeletion.Command.Delete.DomainDelete.Name = domain

	deleteData, err := xml.MarshalIndent(domainDeletion, "", "  ")
	if err != nil {
		return err
	}

	deleteRawResp, err := s.Send(deleteData)
	if err != nil {
		s.logAPIConnectionError(err, "requestID", reqID)
		return err
	}

	var deleteResp epp.APIResult
	if err = xml.Unmarshal(deleteRawResp, &deleteResp); err != nil {
		return err
	}

	if deleteResp.Response.Result.Code != 1000 {
		return errors.New("Request failed: " + deleteResp.Response.Result.Msg)
	}

	return nil
}
